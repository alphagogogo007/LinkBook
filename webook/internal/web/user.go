package web

import (
	"net/http"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

const (
	emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegexPattern = `^.{8,}$`
)

type UserHandler struct{
	emailRexExp *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc *service.UserService

}

func NewUserHandler(svc *service.UserService) *UserHandler{
	return &UserHandler{
		emailRexExp: regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern,regexp.None),
		svc: svc,
	}
}


func (h *UserHandler) RegisterRoutes(server *gin.Engine){


	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
}


func (h *UserHandler) SignUp(ctx *gin.Context){

	type SignUpReq struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
	}
	var req SignUpReq
	if err:= ctx.Bind(&req);err!=nil{
		return 
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil{
		ctx.String(http.StatusOK, "系统错误")
		return 
	}
	if !isEmail{
		ctx.String(http.StatusOK, "非法邮箱格式")
		return 
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil{
		ctx.String(http.StatusOK, "系统错误")
		return 
	}
	if !isPassword{
		ctx.String(http.StatusOK, "密码格式不对，至少需要八位密码")
		return 
	}
	if req.Password!=req.ConfirmPassword{
		ctx.String(http.StatusOK, "两次密码输入不对")
		return 
	}

	err = h.svc.SignUp(ctx,  domain.User{
		Email: req.Email,
		Password: req.Password,
	})

	switch err{
	case nil:
		ctx.String(http.StatusOK, "hello, you are signing up")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}


}


func (h *UserHandler) Login(ctx *gin.Context){


	type LoginReq struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
	}
	var req LoginReq
	if err:= ctx.Bind(&req);err!=nil{
		return 
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err{
	case nil:

		sess := sessions.Default(ctx)
		sess.Set("userId", user.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		err = sess.Save()
		if err!=nil{
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

func (h *UserHandler) Profile(ctx *gin.Context){
	ctx.String(http.StatusOK, "this is your profile")
}

func (h *UserHandler) Edit(ctx *gin.Context){
	
}