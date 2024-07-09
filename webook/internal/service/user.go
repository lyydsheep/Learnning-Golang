package service

import (
	"context"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//需要考虑加密问题
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	//然后就是存数据库
	return svc.repo.Create(ctx, u)
}
