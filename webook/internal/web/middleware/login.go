package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//忽略对于部分路径的请求
		for _, v := range l.paths {
			if v == ctx.Request.URL.Path {
				return
			}
		}
		//检验session
		session := sessions.Default(ctx)
		id := session.Get("UserId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//以Go的方式进行编码解码
		gob.Register(time.Now())
		session.Set("UserId", id)
		session.Options(sessions.Options{
			MaxAge: 60,
		})

		updateTime := session.Get("update_time")
		if updateTime == nil {
			session.Set("update_time", time.Now())
			if err := session.Save(); err != nil {
				panic(err)
			}
			return
		}
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		//对session进行刷新
		if time.Now().Sub(updateTimeVal) > time.Second*10 {
			session.Set("update_time", time.Now())
			if err := session.Save(); err != nil {
				panic(err)
			}
		}
	}
}
