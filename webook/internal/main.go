package main

import (
	"strings"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	login "gitee.com/geekbang/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
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

func initDB() *gorm.DB{

	
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err!=nil{
		panic(err)
	}

	err = dao.InitTables(db)
	if err!=nil{
		panic(err)
	}
	return db

	
}

func initWebServer() *gin.Engine{

	
	server := gin.Default()

	

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))



	login := &login.LoginMiddlewareBuiler{}
	store := cookie.NewStore([]byte("secret"))

	server.Use(	sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}