package main

import (
	"net/http"

	login "gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
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

func useSession(server *gin.Engine) {
	login := &login.LoginMiddlewareBuiler{}
	store, err := redisSession.NewStore(16, "tcp", "localhost:6379", "", []byte("uVCS5zcJSVZjNYoQOJxd9XOYmTUjQ3lP"), []byte("7NcCe8cUJHcaRQa95Xl5isayrYrfijmX"))
	if err != nil {
		panic(err)
	}

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
