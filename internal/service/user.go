package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"swpuforum/internal/domain"
	"swpuforum/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

var (
	ErrInvalidUserOrPassword = errors.New("邮箱或者密码不对")
	ErrUserNotFind           = repository.ErrUserNotFind
)

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) SignUp(ctx context.Context, user *domain.User) error {
	//业务代码 数据进入 数据层的处理
	//加密
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return u.repo.Create(ctx, user)
}

func (u *UserService) Login(ctx context.Context, user *domain.User) (domain.User, error) {
	result, err := u.repo.FindByEmail(ctx, user.Email)
	if errors.Is(err, ErrUserNotFind) {
		return result, ErrInvalidUserOrPassword
	}
	if err != nil {
		return result, err
	}

	//比较哈希
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		//DEBUG
		return result, ErrInvalidUserOrPassword
	}
	return result, nil
}
