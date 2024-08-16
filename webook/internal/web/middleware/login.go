package login

import (
	
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)


type LoginMiddlewareBuiler struct{

}

func (m *LoginMiddlewareBuiler) CheckLogin() gin.HandlerFunc{
	return func(ctx *gin.Context){
		path := ctx.Request.URL.Path
		if path =="/users/signup" || path =="/users/login"{
			return 
		}
		sess := sessions.Default(ctx)
		if sess.Get("userId")==nil{
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return 
		}
	}
}