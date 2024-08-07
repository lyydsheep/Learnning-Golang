package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
)

func InitGin(u *web.UserHandler, f ...gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(f...)
	u.RegisterRoutes(server)
	return server
}
