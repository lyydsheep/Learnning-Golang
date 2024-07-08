package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	dsn := "root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//只在初始化过程panic
	//panic相当于整个goroutine结束
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
	}))

	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)

	u.RegisterRoutes(server)
	server.Run(":8080")

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello world")
	//})
	//server.Run(":8080")
}
