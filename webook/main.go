package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"gitee.com/geekbang/basic-go/webook/config"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/localsms"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	login "gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/webook/pkg/ginx/middleware/ratelimit"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



func main() {

	// db := initDB()
	// redisClient := initRedis()
	
	// server := initWebServer(redisClient)

	// smsSvc := initMemorySMS()
	// codeSvc := initCodeSvc(redisClient, smsSvc)
	// initUserHdl(db, redisClient, codeSvc,server)

	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	})

	server.Run(":8080")

}

func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, 
	codeSvc *service.CodeService,
	
	server *gin.Engine) {

	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(redisClient)

	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)

	hdl := web.NewUserHandler(us, codeSvc)
	hdl.RegisterRoutes(server)
}

func initCodeSvc(redisClient redis.Cmdable, smsSvc sms.Service) *service.CodeService{
	codeCache := cache.NewCodeCache(redisClient)
	codeRepo := repository.NewCodeRepository(codeCache)
	return service.NewCodeService(codeRepo, smsSvc)

}

func initMemorySMS() sms.Service{
	return localsms.NewService()
}


func initDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db

}

func initRedis() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	// 检查 Redis 连接是否成功
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return redisClient
}


func initWebServer(redisClient redis.Cmdable) *gin.Engine {

	server := gin.Default()

	server.Use(cors.New(cors.Config{
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
	}))

	useRedisRatelimiter(server, redisClient)
	useJWT(server)

	return server
}

func useRedisRatelimiter(server *gin.Engine, redisClient redis.Cmdable) {

	redisLimiter := limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)

	server.Use(ratelimit.NewBuilder(redisLimiter).Build())
}

func useJWT(server *gin.Engine) {
	login := &login.LoginJWTMiddlewareBuiler{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &login.LoginMiddlewareBuiler{}
	store, err := redisSession.NewStore(16, "tcp", "localhost:6379", "", []byte("uVCS5zcJSVZjNYoQOJxd9XOYmTUjQ3lP"), []byte("7NcCe8cUJHcaRQa95Xl5isayrYrfijmX"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
