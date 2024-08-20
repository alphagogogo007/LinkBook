package main

import (
	"strings"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
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

	db := initDB()
	server := initWebServer()

	initUserHdl(db, server)

	server.Run(":8080")

}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)

	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db

}

func initWebServer() *gin.Engine {

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


	useRedisRatelimiter(server)
	useJWT(server)

	return server
}

func useRedisRatelimiter(server *gin.Engine){
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
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
