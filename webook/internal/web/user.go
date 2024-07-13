package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service"
	"net/http"
)

type UserHandler struct {
	passwordRegexp *regexp.Regexp
	emailRegexp    *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "[A-Za-z0-9]+([_\\.][A-Za-z0-9]+)*@([A-Za-z0-9\\-]+\\.)+[A-Za-z]{2,6}"
		passwordRegexPattern = "^(?=.*\\d)(?=.*[A-z])[\\da-zA-Z]{1,9}$"
	)
	e := regexp.MustCompile(emailRegexPattern, regexp.None)
	p := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		passwordRegexp: p,
		emailRegexp:    e,
		svc:            svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userRegisterRoutes := server.Group("/users")
	userRegisterRoutes.POST("/signup", u.SignUp)
	//userRegisterRoutes.POST("/login", u.Login)
	userRegisterRoutes.POST("/login", u.LoginJWT)
	userRegisterRoutes.GET("/profile", u.Profile)
	userRegisterRoutes.POST("/edit", u.Edit)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//调用service服务
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "ErrInvalidUserOrPassword")
		return
	}
	//存储session
	session := sessions.Default(ctx)
	session.Set("UserId", user.Id)
	session.Options(sessions.Options{
		MaxAge: 60,
	})
	err = session.Save()

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//调用service服务
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "ErrInvalidUserOrPassword")
		return
	}
	fmt.Println(user)
	//生成jwt-token
	token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString([]byte("DMMzjITr6EpQOOjgUzoRAb440lKd2d3y"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//将jwt-token加入头部
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	//自动返回400 错误码
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

	err = u.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	fmt.Println(req)
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Name      string `json:"name"`
		Birthday  string `json:"birthday"`
		Biography string `json:"biography"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id := sessions.Default(ctx).Get("UserId")
	if id == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	val, ok := id.(int)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err := u.svc.Edit(ctx, domain.User{
		Id:        val,
		Name:      req.Name,
		Birthday:  req.Birthday,
		Biography: req.Biography,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
}

//使用session的Profile函数
//func (u *UserHandler) Profile(ctx *gin.Context) {
//	id := sessions.Default(ctx).Get("UserId")
//	if id == nil {
//		ctx.AbortWithStatus(http.StatusUnauthorized)
//		return
//	}
//	val, ok := id.(int)
//	if !ok {
//		ctx.String(http.StatusOK, "系统错误")
//		return
//	}
//	user, err := u.svc.Profile(ctx, val)
//	if errors.Is(err, service.ErrUserNotFound) {
//		ctx.String(http.StatusOK, "系统错误")
//		return
//	}
//	ctx.String(http.StatusOK, "Name: %s, Birthday: %s Biography: %s", user.Name, user.Birthday, user.Biography)
//}

func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "successful")
}
