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
	db := InitDB()
	server := InitWebServer() //获取虚拟服务器，并通过中间件解决跨域问题

	u := InitUser(db) //初始化有关user的业务准备
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func InitUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func InitWebServer() *gin.Engine {
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
	return server
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"), &gorm.Config{})
	//只在初始化过程panic
	//panic相当于整个goroutine结束
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
