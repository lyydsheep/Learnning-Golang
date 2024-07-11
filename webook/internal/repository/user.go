package repository

import (
	"context"
	"errors"
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

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//通过dao操作数据库
	u, err := repo.ud.FindByEmail(ctx, email)
	if errors.Is(err, ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.ud.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	return repo.ud.Update(ctx, dao.User{
		Id:        u.Id,
		Name:      u.Name,
		Birthday:  u.Birthday,
		Biography: u.Biography,
	})
}

func (repo *UserRepository) FindById(ctx context.Context, id int) (domain.User, error) {
	user, err := repo.ud.FindById(ctx, id)
	if errors.Is(err, ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	return domain.User{
		Name:      user.Name,
		Birthday:  user.Birthday,
		Biography: user.Biography,
	}, err
}
