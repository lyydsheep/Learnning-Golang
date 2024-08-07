package service

import (
	"context"
	"errors"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// var ErrInvalidUserOrPassword = errors.New("账号无效或密码错误")
var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号无效或密码错误")
	ErrUserNotFound          = repository.ErrUserNotFound
)

// basic --- v1
// enhanced --- v2
//

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type BasicUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &BasicUserService{
		repo: repo,
	}
}

func (svc *BasicUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//调用repository服务
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	//密码比较，hash在前，校验在后
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return domain.User{Id: u.Id}, nil
}

func (svc *BasicUserService) SignUp(ctx context.Context, u domain.User) error {
	//需要考虑加密问题
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	//然后就是存数据库
	return svc.repo.Create(ctx, u)
}

func (svc *BasicUserService) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}

func (svc *BasicUserService) Profile(ctx context.Context, id int) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if errors.Is(err, ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	return u, err
}

func (svc *BasicUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, ErrUserNotFound) {
		//存在或错误
		return u, err
	}
	//木有，需要创建
	//慢路径
	err = svc.repo.Create(ctx, domain.User{Phone: phone})
	if err != nil && !errors.Is(err, ErrUserDuplicate) {
		return domain.User{}, err
	}
	// 这里会有主从延迟问题（不懂~_~)
	return svc.repo.FindByPhone(ctx, phone)
}
