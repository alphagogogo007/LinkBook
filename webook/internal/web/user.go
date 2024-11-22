package web

import (
	"log"
	"net/http"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegexPattern = `^.{8,}$`
	JWTKey               = "jYe8vbdGFD7RRnIf8W7KArU2ehZJbbn8"
	bizLogin             = "Login"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {

	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.LoginJWT)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Please input phone number",
		})
	}
	//log.Println(bizLogin, req.Phone)
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	//log.Println(err)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{

			Msg: "Successfully send the code",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Send too many",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})

	}

}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		log.Printf("system error: %v", err)
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "Wrong code",
		})
		return
	}

	// login or create user
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System error",
		})
		return
	}
	h.SetJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "Successfully login",
	})
	return

}

func (h *UserHandler) SignUp(ctx *gin.Context) {

	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		log.Println(err)
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "illegal email name")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "The password format is incorrect; it must be at least eight characters long.")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "The passwords entered do not match")
		return
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "hello, successfully signing up")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "Email conflict, please use a different one.")
	default:
		ctx.String(http.StatusOK, "system error")

	}

}

func (h *UserHandler) Login(ctx *gin.Context) {

	type LoginReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:

		sess := sessions.Default(ctx)
		sess.Set("userId", user.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误: %v", err)
			return
		}
		ctx.String(http.StatusOK, "登录成功")

	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误: %v", err)
	}

}

func (h *UserHandler) SetJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte(JWTKey))
	//log.Printf("login token str: %v", tokenStr)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误: %v", err)
	}
	ctx.Header("x-jwt-token", tokenStr)
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {

	type LoginReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.SetJWTToken(ctx, user.Id)
		//log.Println("登录成功， tokenStr")
		ctx.String(http.StatusOK, "登录成功")

	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误: %v", err)
	}

}

func (h *UserHandler) Profile(ctx *gin.Context) {

	// userId, err := h.svc.GetUserIdFromSession(ctx)
	// if err != nil {
	// 	ctx.AbortWithStatus(http.StatusUnauthorized)
	// 	ctx.String(http.StatusUnauthorized, "系统错误: %v", err)
	// 	return
	// }

	us, ok := ctx.MustGet("user").(UserClaims)
	if !ok {
		log.Println("Type assertion to UserClaims failed")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := us.Uid

	u, err := h.svc.FindById(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误: %v", err)
		return
	}

	//json大小写要匹配！！！前端是Nickname, 所以你也是json:Nickname
	type User struct {
		Nickname string `json:"Nickname"`
		Email    string `json:"Email"`
		AboutMe  string `json:"AboutMe"`
		Birthday string `json:"Birthday"`
	}

	frontUserProfile := User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	}

	ctx.JSON(http.StatusOK, frontUserProfile)
}

func (h *UserHandler) Edit(ctx *gin.Context) {

	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutme"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	userId, err := h.svc.GetUserIdFromSession(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.String(http.StatusUnauthorized, "系统错误: %v", err)
		return
	}

	//println(userId, "edit userId")

	// check birthday
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}

	if err := h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       userId,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	}); err != nil {

		ctx.String(http.StatusInternalServerError, "edit profile error:%v", err)
		return
	}
	//println("success")
	ctx.String(http.StatusOK, "Edit successful")
}
