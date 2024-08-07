package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web/middleware"
	"github.com/lyydsheep/Learnning-Golang/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitMiddleware(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//解决CORS中间件
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			ExposeHeaders: []string{"x-jwt-token"},
		}),
		//限流中间件
		ratelimit.NewBuilder(redisClient, time.Minute, 100).Build(),
		//登录状态校验中间件
		(middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePath("/users/signup")).
			IgnorePath("/users/login").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			Build(),
	}
}
