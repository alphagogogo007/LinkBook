package ioc

import (
	//"os"

	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service{
	// appID, ok := os.LookupEnv("WECHAT_APP_ID")
	// if !ok{
	// 	panic("cannot find global environment variable wechat")
	// }
	appID := "wxbdc5610cc59c1631"
	appSecret := "abcdefg"
	return wechat.NewService(appID, appSecret)
}