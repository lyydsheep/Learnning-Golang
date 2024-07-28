package repository

import (
	"context"
	"errors"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/cache"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	ud *dao.UserDAO
	uc *cache.UserCache
}

func NewUserRepository(ud *dao.UserDAO, uc *cache.UserCache) *UserRepository {
	return &UserRepository{
		ud: ud,
		uc: uc,
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
	u, err := repo.uc.Get(ctx, id)
	//缓存中有数据
	if err == nil {
		return u, err
	}
	//缓存中没有数据 ---> 查数据库
	//缓存炸了 ---> 直接查数据库，有很大的风险，最好做到数据库的限流
	ue, err := repo.ud.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	//在数据库中查到了数据，设置缓存
	u = domain.User{
		Id:        ue.Id,
		Name:      ue.Name,
		Birthday:  ue.Birthday,
		Biography: ue.Biography,
	}
	err = repo.uc.Set(ctx, u)
	if err != nil {
		//日志监控
	}
	return u, err
}
