//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/cache"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
	"github.com/lyydsheep/Learnning-Golang/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 初始化第三方服务
		ioc.InitDB, ioc.InitRedis, ioc.InitSMS,
		// 初始化DAO和cache
		dao.NewUserDAO, cache.NewUserCache, cache.NewCodeRedis,
		// 初始化repository
		repository.NewUserRepository, repository.NewCodeRepository,
		// 初始化service
		service.NewUserService, service.NewCodeService,
		// 初始化web（UserHandler）
		web.NewUserHandler,
		// 初始化中间件
		ioc.InitMiddleware,
		// 初始化gin----使用中间件、进行路由注册
		ioc.InitGin,
	)
	return new(gin.Engine)
}
