package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/service"
	"net/http"
	"time"
)

// 不是controller，是handler
// 代理函数，相当于helper

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	// 用户名必须是大小写字母和下划线组合，长度为12以内
	nickNameRegexPattern = `^[a-zA-Z0-9_]{1,12}$`
	bizLogin             = "login"
)

type UserHandler struct {
	svc            service.UserService
	codeSvc        service.CodeService
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	nickNameRegExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:            svc,
		codeSvc:        codeSvc,
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		nickNameRegExp: regexp.MustCompile(nickNameRegexPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/signup", h.SignUp)
		ug.POST("/login", h.LoginJWT)
		ug.POST("/edit", h.Edit)
		ug.GET("/profile", h.Profile)
		ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
		ug.POST("/login_sms", h.LoginSMS)
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	// 解析前端发来的数据
	// 拿到用户的username和pwd
	// 注册之
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，换一个邮箱")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	// refresh expiry time for login
	// change other stuff
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	userClaim, ok := ctx.MustGet("user").(UserClaim)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 用户输入非法用户名
	isNickName, err := h.nickNameRegExp.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusOK, "非法用户名")
		return
	}
	errorResp := make(map[string]string)
	if !isNickName {
		errorResp["nickname"] = "用户名必须是大小写字母和下划线组合，长度为12以内"
	}
	// 用户输入不对
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		errorResp["birthday"] = "生日格式错误"
	}
	// 用户输入超过140个字个人简介
	if len(req.AboutMe) > 140 {
		errorResp["about_me"] = "个人简介限制在140个字以内"
	}

	if len(errorResp) > 0 {
		ctx.JSON(http.StatusOK, errorResp)
		return
	}

	err = h.svc.UpdateNonPII(ctx, domain.User{
		Id:       userClaim.Id,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	userClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaim)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	// change information
	ctx.String(http.StatusOK, "更新成功")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	userClaim, ok := ctx.MustGet("user").(UserClaim)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// find by user's id
	u, err := h.svc.FindById(ctx, userClaim.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	// return a selected user
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}

	userClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaim)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	// change information
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		userClaim := UserClaim{
			Id:        user.Id,
			UserAgent: ctx.GetHeader("User-Agent"),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		// issue a token
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaim)
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// 验证有没有这条记录
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 记录登录态
		// session版本
		session := sessions.Default(ctx)
		session.Set("userId", user.Id)
		session.Options(sessions.Options{
			MaxAge: 60,
		})
		err = session.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	// 生成并发送code，并且设置code
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	case service.ErrTooManyCodeSend:
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "短信发送太频繁，请稍后再试"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		// zap log here
	}
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}
	if req.Code == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入验证码"})
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		// 要打印日志
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码不正确"})
		return
	}
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "系统错误"})
		return
	}
	ssid := uuid.New().String()
	// set JWT Token
	err = h.setJWTToken(ctx, ssid, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "登录成功"})
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaim struct {
	jwt.RegisteredClaims
	Id        int64
	UserAgent string
	Ssid      string
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, ssid string, uid int64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaim{
		Id:        uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
