package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/dao"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web/middleware"
	"github.com/lyydsheep/Learnning-Golang/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := InitDB()
	server := InitWebServer() //获取虚拟服务器，并通过中间件解决跨域问题

	u := InitUser(db)        //初始化有关user的业务准备
	u.RegisterRoutes(server) //对user有关的业务进行路由注册

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
	//解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
		ExposeHeaders: []string{"x-jwt-token"},
	}))

	//初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	//使用gin中间件进行限流
	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	////创建session
	//store, err := redis.NewStore(32, "tcp", "localhost:6379", "",
	//	[]byte("fD6TyDBMbRsRYZW3PWI6y4r5oeLJv2x38kSXNHgn6raksxXuIzheW0Bgd6BiVrv0"),
	//	[]byte("xZuTExq1NQFqFNvoMykWrmhtvzOP4rM8"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//server.Use(sessions.Sessions("mySession", store))

	server.Use((middleware.NewLoginJWTMiddlewareBuilder().IgnorePath("/users/signup")).
		IgnorePath("/users/login").Build())

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
