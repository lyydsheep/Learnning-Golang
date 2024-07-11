package repository

import (
	"context"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	ud *dao.UserDAO
}

func NewUserRepository(ud *dao.UserDAO) *UserRepository {
	return &UserRepository{
		ud: ud,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, id int) (domain.User, error) {
	u, err := repo.ud.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Name:      u.Name,
		Birthday:  u.Birthday,
		Biography: u.Biography,
	}, nil
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	return repo.ud.Update(ctx, dao.User{
		Id:        u.Id,
		Name:      u.Name,
		Biography: u.Biography,
		Birthday:  u.Birthday,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.ud.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.ud.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
