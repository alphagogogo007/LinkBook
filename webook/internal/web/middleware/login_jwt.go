package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuiler struct {
}

func (m *LoginJWTMiddlewareBuiler) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			log.Println("authcode empty")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//log.Println("authcode :", authCode)
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte(web.JWTKey), nil
		})
		if err != nil {
			log.Println("parse token error")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token==nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent!=ctx.GetHeader("User-Agent"){
			// Need to track events here
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		//log.Println("expire time", expireTime.Time)
		// if expireTime.Before(time.Now()){
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }
		//log.Println("token str:", tokenStr)
		if expireTime.Sub(time.Now()) < time.Minute*3 {
			//刷新 jwt token
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute*30))
			//token中有一个claim指针，所以修改uc.Expire就能直接对token生成tokenstr产生影响
			tokenStr, err = token.SignedString([]byte(web.JWTKey))
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println("refresh error", err)
			}

		}
		// uc里面有uid
		ctx.Set("user", uc)
	}
}
