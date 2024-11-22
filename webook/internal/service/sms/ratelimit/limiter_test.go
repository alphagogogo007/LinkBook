package ratelimit

import (
	"context"
	"errors"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	type fields struct {
		svc     sms.Service
		limiter limiter.Limiter
		key     string
	}

	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "no limit",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc,l
			},
			wantErr: nil,
		},
		{
			name: "has limit",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				
				return svc,l
			},
			wantErr: errLimited,
		},
		{
			name: "system error",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("redis limiter error"))
				
				return svc,l
			},
			wantErr: errors.New("redis limiter error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tt.mock(ctrl)
			svc := NewRateLimitSMSService(smsSvc, l)

			err := svc.Send(context.Background(), "abc", []string{"123"}, "123456")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
