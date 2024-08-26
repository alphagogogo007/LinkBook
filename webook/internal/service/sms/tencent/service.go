package tencent

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)


// Client is an interface that wraps the SendSms method.
type Client interface {
	SendSms(request *sms.SendSmsRequest) (*sms.SendSmsResponse, error)
}

// 这里的client本来是*sms.Client
type Service struct {
	client   Client
	appId    *string
	signName *string
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		return fmt.Errorf("an API error has returned: %v", err)

	}

	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "ok" {
			return fmt.Errorf("sms failed. code:%s, msg:%s", *status.Code, *status.Message)

		}
	}
	return nil

}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

func NewService(client Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}
