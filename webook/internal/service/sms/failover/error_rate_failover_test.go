package failover

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestErrorRateFailoverSMSService_Send(t *testing.T) {
	tests := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) []sms.Service
		threshold  float64
		windowSize time.Duration
		setup      func(svc *ErrorRateFailoverSMSService)
		wantIdx    int32
		wantErr    error
	}{
		{
			name: "no error, no failover",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			threshold:  0.1,
			windowSize: time.Minute,
			setup: func(svc *ErrorRateFailoverSMSService) {
				svc.requestCh <- time.Now()
			},
			wantIdx: 0,
			wantErr: nil,
		},
		{
			name: "error rate below threshold, no failover",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				//svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil) 
				// 如果不调用，那么就不要写，不然会报错
				return []sms.Service{svc0, svc1}
			},
			threshold:  0.4,
			windowSize: time.Minute,
			setup: func(svc *ErrorRateFailoverSMSService) {
				for i := 0; i < 10; i++ {
					svc.requestCh <- time.Now()
				}
				for i := 0; i < 2; i++ {
					svc.errorCh <- time.Now()
				}
			},//这里有问题，这里应该超过threshold了
			wantIdx: 0,
			wantErr: nil,
		},
		{
			name: "error rate exceeds threshold, trigger failover",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			threshold:  0.1,
			windowSize: time.Minute,
			setup: func(svc *ErrorRateFailoverSMSService) {
				for i := 0; i < 10; i++ {
					svc.requestCh <- time.Now()
				}
				for i := 0; i < 2; i++ {
					svc.errorCh <- time.Now()
				}
			},
			wantIdx: 1,
			wantErr: nil,
		},
		{
			name: "error rate exceeds threshold, new service fails",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("service failed"))
				return []sms.Service{svc0, svc1}
			},
			threshold:  0.1,
			windowSize: time.Minute,
			setup: func(svc *ErrorRateFailoverSMSService) {
				for i := 0; i < 10; i++ {
					svc.requestCh <- time.Now()
				}
				for i := 0; i < 2; i++ {
					svc.errorCh <- time.Now()
				}
			},
			wantIdx: 1,
			wantErr: errors.New("service failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock services
			svcs := tt.mock(ctrl)

			// Initialize service
			service := NewErrorRateFailoverSMSService(svcs, tt.threshold, tt.windowSize)
			defer service.Stop() // Ensure cleanup after test

			// Setup preconditions
			tt.setup(service)

			// Call the Send method
			err := service.Send(context.Background(), "tplID", []string{"arg1"}, "123456789")

			// Assertions
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantIdx, service.idx)
		})
	}
}
