package repository

import (
	"context"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
)

type UserRepository struct {
	ud *dao.UserDAO
}

func NewUserRepository(ud *dao.UserDAO) *UserRepository {
	return &UserRepository{
		ud: ud,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.ud.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
