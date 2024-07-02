package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	passwordRegexp *regexp.Regexp
	emailRegexp    *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = "[A-Za-z0-9]+([_\\.][A-Za-z0-9]+)*@([A-Za-z0-9\\-]+\\.)+[A-Za-z]{2,6}"
		passwordRegexPattern = "^(?=.*\\d)(?=.*[A-z])[\\da-zA-Z]{1,9}$"
	)
	p := regexp.MustCompile(passwordRegexPattern, regexp.None)
	e := regexp.MustCompile(emailRegexPattern, regexp.None)
	return &UserHandler{
		passwordRegexp: p,
		emailRegexp:    e,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userRegisterRouters := server.Group("/users")
	userRegisterRouters.POST("/signup", u.SignUp)
}

func (u *UserHandler) Login(ctx *gin.Context) {}
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	//校验邮箱格式
	ok, err := u.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}
	//校验密码格式
	ok, err = u.passwordRegexp.MatchString(req.Password)
	if err != nil {
		//记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "只能由字母、数字组成，1-9位")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {}

func (u *UserHandler) Profile(ctx *gin.Context) {}
