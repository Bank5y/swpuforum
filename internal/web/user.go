package web

import (
	"errors"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"swpuforum/internal/domain"
	"swpuforum/internal/repository/dao"
	"swpuforum/internal/service"
)

type UserHandler struct {
	svc         *service.UserService
	EmailExp    *regexp2.Regexp
	PasswordExp *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		//email regex
		emailRegexPattern = `^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
		//password regex
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	EmailExp := regexp2.MustCompile(emailRegexPattern, regexp2.None)
	PasswordExp := regexp2.MustCompile(passwordRegexPattern, regexp2.None)
	return &UserHandler{
		svc:         svc,
		EmailExp:    EmailExp,
		PasswordExp: PasswordExp,
	}
}

func (u *UserHandler) RegisterRouter(server *gin.Engine) {
	usersGroup := server.Group("/users")
	usersGroup.POST("/signup", u.SignUp)
	usersGroup.POST("/login", u.login)
	usersGroup.GET("/profile", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你看到了。。。")
	})
}

func (u *UserHandler) SignUp(ctx *gin.Context) {

	//请求体绑定
	type SignUpReq struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ConfirmedPassword string `json:"confirmedPassword"`
	}
	var signUpReq SignUpReq
	err := ctx.Bind(&signUpReq)

	if err != nil {
		return
	}
	//邮箱正则验证
	emailMatch, err := u.EmailExp.MatchString(signUpReq.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !emailMatch {
		ctx.String(http.StatusOK, "邮箱有误")
		return
	}
	//密码正则验证
	passwordMatch, err := u.PasswordExp.MatchString(signUpReq.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !passwordMatch {
		ctx.String(http.StatusOK, "密码格式错误！")
		return
	}
	//密码确认验证
	if signUpReq.Password != signUpReq.ConfirmedPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致！")
	}

	err = u.svc.SignUp(ctx.Request.Context(), &domain.User{
		Email:    signUpReq.Email,
		Password: signUpReq.Password,
	})
	if errors.Is(err, dao.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "注册成功")

}

func (u *UserHandler) login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	result, err := u.svc.Login(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//设置session
	sess := sessions.Default(ctx)
	sess.Set("LoginSess", result.Email)
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}
