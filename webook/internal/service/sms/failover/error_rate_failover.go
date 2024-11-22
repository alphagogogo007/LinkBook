package failover

import (
	"context"
	"sync/atomic"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)

type ErrorRateFailoverSMSService struct {
	svcs       []sms.Service
	idx        int32
	threshold  float64
	windowSize time.Duration
	requestCh  chan time.Time
	errorCh    chan time.Time
	stopCh     chan struct{} //what is this?

}

func NewErrorRateFailoverSMSService(svcs []sms.Service, threshold float64, windowSize time.Duration) *ErrorRateFailoverSMSService {
	service := &ErrorRateFailoverSMSService{
		svcs:       svcs,
		threshold:  threshold,
		windowSize: windowSize,
		requestCh:  make(chan time.Time, 100000),
		errorCh:    make(chan time.Time, 100000),
		stopCh:     make(chan struct{}),
	}
	go service.cleanupExpiredRequest()
	return service
}

func (e *ErrorRateFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {


	idx := e.idx
	//println(e.calculateErrorate())
	if e.calculateErrorate() > e.threshold {
		newIdx := (idx + 1) % int32(len(e.svcs))
		if atomic.CompareAndSwapInt32(&e.idx, idx, newIdx) {
			e.clearChannels()
		}
		idx = newIdx
		//println(idx)
	}

	svc := e.svcs[idx]

	err := svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		e.errorCh <- time.Now()
	}
	e.requestCh <- time.Now()

	return err
}

func (e *ErrorRateFailoverSMSService) cleanupExpiredRequest(){
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			cutoff := time.Now().Add(-e.windowSize)
			e.cleanupChannel(e.requestCh, cutoff)
			e.cleanupChannel(e.errorCh, cutoff)

		case <-e.stopCh: // what is this stopch??
			return
		}
	}
}

func (e *ErrorRateFailoverSMSService) cleanupChannel(ch chan time.Time, cutoff time.Time){
	for {
		select{
		case t:= <-ch:
			if t.After(cutoff){
				ch<-t //如果数据量大，每秒放回去1个不影响总量和error统计
					  //如果数据量小，很快就会比如1秒5个，不到1分钟就会被清除出去
				return
			}
		default:
			return
		
		}
	
	}
}

func (e *ErrorRateFailoverSMSService) calculateErrorate()  float64{
	requests := len(e.requestCh)
	if requests==0{
		return 0
	}
	errors := len(e.errorCh)
	return float64(errors)/float64(requests)
}


func (e *ErrorRateFailoverSMSService) clearChannels(){
	for len(e.requestCh)>0{
		<-e.requestCh
	}
	for len(e.errorCh)>0{
		<-e.errorCh
	}
}

func (e *ErrorRateFailoverSMSService) Stop(){
	close(e.stopCh) //还是没有搞清楚e.stopCh有什么用
}

