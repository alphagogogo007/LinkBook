package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"


	"gitee.com/geekbang/basic-go/webook/internal/domain"
	uuid "github.com/lithammer/shortuuid/v4"

)

type Service interface{
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)

}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

//var redirectURL = url.PathEscape("https://passport.yhd.com/wechat/callback.do")

type service struct{
	appID string
	appSecret string
	client *http.Client
}

func NewService(appID string, appSecret string) Service{
	return &service{
		appID: appID,
		appSecret: appSecret,
		client: http.DefaultClient,
	}
}

func (s *service) VerifyCode(ctx context.Context, 
	code string) (domain.WechatInfo, error){

		acessTokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code ) 
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, acessTokenURL, nil)
		if err!=nil{
			return domain.WechatInfo{}, err
		}
		httpResp, err := s.client.Do(req)
		if err!=nil{
			return domain.WechatInfo{}, err
		}

		var res Result
		err = json.NewDecoder(httpResp.Body).Decode(&res)
		if err!=nil{
			return domain.WechatInfo{}, err
		}
		if res.ErrCode!=0{
			return domain.WechatInfo{}, fmt.Errorf("wechat api cal failure, error code:%d, error msg:%s",res.ErrCode, res.ErrMsg)
		}
	
		return domain.WechatInfo{
			UnionId: res.UnionId,
			OpenId: res.OpenId,
		}, nil
}

func (s *service) AuthURL(ctx context.Context) (string, error){
	
	state := uuid.New()
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil

}

type Result struct{
	AccessToken	string `json:"access_token"`
	ExpiresIn	int64 `json:"expires_in"`
	RefreshToken	string `json:"refresh_token"`
	OpenId	string `json:"openid"`
	Scope	string `json:"scope"`
	UnionId	string 	`json:"union_id"`
	ErrCode int64 `json:"errcode"`
	ErrMsg string  `json:"errmsg"`

}