package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
}

const JWTKey = "jYe8vbdGFD7RRnIf8W7KArU2ehZJbbn8"

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}


func (h *jwtHandler) SetJWTToken(ctx *gin.Context, uid int64) {
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
