//go:build wireinject
package startup

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

//这些代码是干嘛的，明明已经有wire了？？？用于integration test的
func InitWebServer() *gin.Engine {

	wire.Build(
		//context
		//ioc.InitContext,

		//third party
		ioc.InitDB,
		InitRedis,
		//ioc.InitBigCache,

		//dao
		dao.NewUserDao,

		//cache
		cache.NewRedisCodeCache, 
		//cache.NewBigCacheCodeCache,
		cache.NewRedisUserCache,
		

		//repository
		repository.NewCodeRepository,
		repository.NewUserRepository,

		//service
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,

		//handler
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,

		ioc.NewLimiter,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
