//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"week2/wire/repository"
	"week2/wire/repository/dao"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return new(repository.UserRepository)
}
