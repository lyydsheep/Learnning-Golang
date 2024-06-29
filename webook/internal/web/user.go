package web

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/users/login", u.Login)
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) SignUp(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
