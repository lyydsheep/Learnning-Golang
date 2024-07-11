package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, v := range l.paths {
			if v == ctx.Request.URL.Path {
				return
			}
		}
		session := sessions.Default(ctx)
		id := session.Get("userID")
		if id == nil {
			fmt.Println("this is LoginMiddleware")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
