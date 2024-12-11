package ioc

import (
	"context"
	"strings"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/web"
	login "gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middleware/ratelimit"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitContext() context.Context{
	return context.Background()
}

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server

}

func InitGinMiddlewares(redisLimiter limiter.Limiter) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		ratelimit.NewBuilder(redisLimiter).Build(),
		(&login.LoginJWTMiddlewareBuiler{}).CheckLogin(),
		
	}
}

func NewLimiter(redisClient redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)
}
