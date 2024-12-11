package web

import (
	"net/http"

	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
)


type OAuth2WechatHandler struct{
	jwtHandler
	svc wechat.Service
	userSvc service.UserService
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler{
	return &OAuth2WechatHandler{
		svc: svc,
		userSvc: userSvc,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine){
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context){
	val, err := o.svc.AuthURL(ctx)
	if err!=nil{
		ctx.JSON(http.StatusOK, Result{
			Msg: "construct ulr error",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})

}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context){
	
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