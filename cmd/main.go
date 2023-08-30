package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"swpuforum/internal/repository"
	"swpuforum/internal/repository/dao"
	"swpuforum/internal/service"
	"swpuforum/internal/web"
	"swpuforum/internal/web/middleware"
)

func main() {
	server := initWebServer()
	db := initDB()

	//注册User业务逻辑
	handler := initUserHandler(db)
	handler.RegisterRouter(server)

	server.Run(":8080")
}

// User业务
func initUserHandler(db *gorm.DB) *web.UserHandler {
	userDAO := dao.NewUserDAO(db)
	repo := repository.NewUserRepo(userDAO)
	userService := service.NewUserService(repo)
	handler := web.NewUserHandler(userService)
	return handler
}

// 初始化中间件
func initWebServer() *gin.Engine {
	server := gin.Default()
	//跨域
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost"},
		AllowCredentials: true,
	}))

	//session处理
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mySessions", store))

	//验证登录状态
	server.Use(
		middleware.NewLoginMiddlewareBuilder().
			Ignore("/users/login").
			Ignore("/users/signup").
			Build())

	return server
}

// 初始化数据库
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/swpu"))
	if err != nil {
		//panic相当于整个goroutine结束
		//整个goroutine结束
		panic(err)
	}
	//建表初始化
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
