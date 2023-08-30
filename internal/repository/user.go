package repository

import (
	"context"
	"swpuforum/internal/domain"
	"swpuforum/internal/repository/dao"
)

type UserRepo struct {
	userDAO *dao.UserDAO
}

var (
	ErrUserNotFind = dao.ErrUserNotFind
)

func NewUserRepo(userDAO *dao.UserDAO) *UserRepo {
	return &UserRepo{userDAO: userDAO}
}

func (repo *UserRepo) Create(ctx context.Context, user *domain.User) error {
	return repo.userDAO.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, err
}
