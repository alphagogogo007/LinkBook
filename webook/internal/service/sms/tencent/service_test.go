package tencent

import (
	"context"
	"testing"
	
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct{}

// SendSms is a mock implementation of the SendSms method.
func (m *MockClient) SendSms(request *sms.SendSmsRequest) (*sms.SendSmsResponse, error) {
	// Return a mock response
	return &sms.SendSmsResponse{
		Response: &sms.SendSmsResponseParams{
			SendStatusSet: []*sms.SendStatus{
				{
					Code:    ekit.ToPtr("ok"),
					Message: ekit.ToPtr("Success"),
				},
			},
		},
	}, nil
}

//TO DO 
//这个需要手动跑，也就是你需要在本地搞好这些环境变量
func TestSender(t *testing.T) {
	mockClient := &MockClient{}
	s := NewService(mockClient, "1400842696", "妙影科技")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "发送验证码",
			tplId:  "1877556",
			params: []string{"123456"},
			// 改成你的手机号码
			numbers: []string{"+8613711112222"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
