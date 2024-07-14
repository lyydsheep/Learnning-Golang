package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/web"
	"net/http"
	"strings"
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
		//将JWT-token中所蕴涵的userId赋给claims
		uc := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			//使用加密的key进行解密
			return []byte("6dGChSIkiB7LRnrpSiYgRe1gtbPdbXit"), nil
		})
		//没登录
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//没登录
		if token == nil || uc.UserId == 0 || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println(uc.UserId)
		//将userId存入上下文中，方便后续获取数据
		ctx.Set("userClaims", uc)
	}
}
