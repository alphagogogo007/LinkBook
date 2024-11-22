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

func TestTimeoutFailoverSMSService_Send(t *testing.T) {

	tests := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.Service
		idx       int32
		cnt       int32
		threshold int32
		wantIdx int32
		wantCnt int32
		wantErr   error
	}{
		// TODO: Add test cases.
		{
			name: "no swap",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			idx: 0,
			cnt: 12,
			threshold: 15,
			wantErr: nil,
			wantIdx: 0,
			wantCnt: 0,
		},
		{
			name: "trigger swap",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
	
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx: 0,
			cnt: 15,
			threshold: 15,
			wantErr: nil,
			wantIdx: 1,
			wantCnt: 0,
		},
		{
			name: "trigger swap, then failed",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
	
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed"))
				return []sms.Service{svc0, svc1}
			},
			idx: 1,
			cnt: 15,
			threshold: 15,
			wantErr: errors.New("failed"),
			wantIdx: 0,
			wantCnt: 0,
		},
		{
			name: "trigger swap, then exceed time",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
	
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			idx: 1,
			cnt: 15,
			threshold: 15,
			wantErr: context.DeadlineExceeded,
			wantIdx: 0,
			wantCnt: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSMSService(tt.mock(ctrl), tt.threshold)
			svc.cnt = tt.cnt
			svc.idx = tt.idx
			err := svc.Send(context.Background(), "123", []string{"1234"}, "12345")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantCnt, svc.cnt)
			assert.Equal(t, tt.wantIdx, svc.idx)
			
		})
	}
}
