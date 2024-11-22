package failover

import (
	"context"
	"errors"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFailoverSMSService_Send(t *testing.T) {

	tests := []struct {
		name    string
		mock func(ctrl *gomock.Controller) []sms.Service


		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "sucess once",
			mock : func(ctrl *gomock.Controller) []sms.Service{
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil )
				return []sms.Service{svc0}
			},
			wantErr: nil,
		},
		{
			name: "sucess at second",
			mock : func(ctrl *gomock.Controller) []sms.Service{
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send failed"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil )
				return []sms.Service{svc0, svc1}
			},
			wantErr: nil,
		},
		{
			name: "all failed",
			mock : func(ctrl *gomock.Controller) []sms.Service{
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send failed"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send failed") )
				return []sms.Service{svc0, svc1}
			},
			wantErr: errors.New("all service failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailoverSMSService(tt.mock(ctrl))
			
			err := svc.Send(context.Background(), "123", []string{"1234"}, "12345")
			assert.Equal(t, tt.wantErr, err)
			})
	}
}
