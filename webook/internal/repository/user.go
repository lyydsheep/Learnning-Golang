package repository

import (
	"context"
	"database/sql"
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
	return repo.entityToDomain(u), err
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.ud.Insert(ctx, repo.domainToEntity(u))
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	return repo.ud.Update(ctx, repo.domainToEntity(u))
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
	u = repo.entityToDomain(ue)
	err = repo.uc.Set(ctx, u)
	if err != nil {
		//日志监控
	}
	return u, err
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	ue, err := repo.ud.FindByPhone(ctx, phone)
	return repo.entityToDomain(ue), err
}

func (repo *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Name:      u.Name,
		Birthday:  u.Birthday,
		Biography: u.Biography,
	}
}

func (repo *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:        u.Id,
		Email:     u.Email.String,
		Password:  u.Password,
		Phone:     u.Phone.String,
		Name:      u.Name,
		Birthday:  u.Birthday,
		Biography: u.Biography,
	}
}
