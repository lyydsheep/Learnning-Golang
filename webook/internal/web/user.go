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
	"time"
)

type UserHandler struct {
	passwordRegexp *regexp.Regexp
	emailRegexp    *regexp.Regexp
	phoneRegexp    *regexp.Regexp
	codeRegexp     *regexp.Regexp
	svc            *service.UserService
	codeSvc        *service.CodeService
}

const biz = "login"

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "[A-Za-z0-9]+([_\\.][A-Za-z0-9]+)*@([A-Za-z0-9\\-]+\\.)+[A-Za-z]{2,6}"
		passwordRegexPattern = "^(?=.*\\d)(?=.*[A-z])[\\da-zA-Z]{1,9}$"
		phoneRegexpPattern   = "^[1]{1}[0-9]{10}$"
		codeRegexpPattern    = "^[0-9]{6}$"
	)
	e := regexp.MustCompile(emailRegexPattern, regexp.None)
	p := regexp.MustCompile(passwordRegexPattern, regexp.None)
	ph := regexp.MustCompile(phoneRegexpPattern, regexp.None)
	c := regexp.MustCompile(codeRegexpPattern, regexp.None)
	return &UserHandler{
		passwordRegexp: p,
		emailRegexp:    e,
		phoneRegexp:    ph,
		codeRegexp:     c,
		svc:            svc,
		codeSvc:        codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	userRegisterRoutes := server.Group("/users")
	userRegisterRoutes.POST("/signup", u.SignUp)
	//userRegisterRoutes.POST("/login", u.Login)
	userRegisterRoutes.POST("/login", u.LoginJWT)
	userRegisterRoutes.GET("/profile", u.Profile)
	userRegisterRoutes.POST("/edit", u.Edit)
	userRegisterRoutes.POST("/login_sms/code/send", u.LoginSendCode)
	userRegisterRoutes.POST("/login_sms", u.LoginSms)
}

func (u *UserHandler) LoginSendCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ok, err := u.phoneRegexp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "手机号格式有误"})
		return
	}
	err = u.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case errors.Is(err, nil):
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	case errors.Is(err, service.ErrTooFrequent):
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "发送过于频繁"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
}

func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ok, err := u.phoneRegexp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "手机号格式有误"})
		return
	}
	ok, err = u.codeRegexp.MatchString(req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码格式不对"})
		return
	}
	err = u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrNotMatch) {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码有误"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	err = u.setJWT(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "登录/注册成功"})
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

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int
	UserAgent string
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

	if err = u.setJWT(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) setJWT(ctx *gin.Context, uid int) error {
	//将用户id赋值于claims
	uc := UserClaims{
		//添加过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		//添加UserId
		UserId: uid,
		//添加UserAgent
		UserAgent: ctx.Request.UserAgent(),
	}
	//通过claims将UserId、UserAgent进行加密形成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte("6dGChSIkiB7LRnrpSiYgRe1gtbPdbXit"))
	if err != nil {
		return err
	}

	//将jwt-token加入头部
	ctx.Header("x-jwt-token", tokenStr)
	return nil
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
	val, ok := ctx.Get("userClaims")
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	uc, ok := val.(UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	uid := uc.UserId
	user, err := u.svc.Profile(ctx, uid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
	ctx.JSON(http.StatusOK, Result{Msg: fmt.Sprintf("id: %d, name: %s, phone: %s", user.Id, user.Name, user.Phone)})
}
