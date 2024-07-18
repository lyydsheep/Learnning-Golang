package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//忽略对于部分路径的请求
		for _, v := range l.paths {
			if v == ctx.Request.URL.Path {
				return
			}
		}
		//取出Authorization字段的value
		segs := strings.Split(ctx.GetHeader("Authorization"), " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//从Authorization中取出JWT-token
		tokenStr := segs[1]
		//将JWT-token中所蕴涵的信息赋给claims
		claims := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			//使用加密的key进行解密
			return []byte("6dGChSIkiB7LRnrpSiYgRe1gtbPdbXit"), nil
		})
		//没登录
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//没登录
		if token == nil || claims.UserId == 0 || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//每十秒生成一个新的token
		if claims.ExpiresAt.Sub(time.Now()) < time.Second*50 {
			//延续时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("6dGChSIkiB7LRnrpSiYgRe1gtbPdbXit"))
			ctx.Header("x-jwt-token", tokenStr)
		}
		fmt.Println(claims.UserId)
		//将userId存入上下文中，方便后续获取数据
		ctx.Set("userClaims", claims)
	}
}
