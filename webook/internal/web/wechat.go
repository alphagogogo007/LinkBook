package web

import (
	"fmt"
	"net/http"

	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
)


type OAuth2WechatHandler struct{
	jwtHandler
	svc wechat.Service
	userSvc service.UserService
	key []byte
	stateCookieName string
}

type StateClaims struct{
	jwt.RegisteredClaims
	State string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler{
	return &OAuth2WechatHandler{
		svc: svc,
		userSvc: userSvc,
		key: []byte("jYe8vbdGFD7RRnIf8W7KArU2ehZJbbn8"),
		stateCookieName: "jwt-state",
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine){
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context){

	state := uuid.New()
	// 为什么要传到这个svc里面
	val, err := o.svc.AuthURL(ctx, state)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "construct ulr error",
			Code: 5,
		})
		return
	}

	//为什么这里state又被传了一遍
	err = o.SetStateCookie(ctx, state)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "service error",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})

}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context){

	err:= o.VerifyState(ctx)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "illegal request",
			Code: 4,
		})
		return
	}
	
	// code为什么能从ctx拿出来？
	code := ctx.Query("code")
	//state := ctx.Query("state")

	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "authorization code error",
			Code: 4,
		})
		return
	}

	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "system error",
			Code: 5,
		})
		return
	}
	o.SetJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
	return 
}


func (o *OAuth2WechatHandler) VerifyState(ctx *gin.Context) error{
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err!=nil{
		return fmt.Errorf("cannot get cookie %w", err)

	}

	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(t *jwt.Token) (interface{}, error) {
		return o.key, nil
	} )
	if err!=nil{
		return fmt.Errorf("parse token error, %w", err)
	}
	if state!= sc.State{
		return fmt.Errorf("state not match")
	}
	return nil


}
 

// 在哪里使用？在auth2url中使用
func (o *OAuth2WechatHandler) SetStateCookie(ctx *gin.Context, state string) error{

	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err :=  token.SignedString(o.key)
	if err!=nil{
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}