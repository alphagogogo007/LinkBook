package failover

import (
	"context"
	"sync/atomic"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)

type AsyncFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	cnt  int32
	curTime int64
	threshold int32
	timeDuration time.Duration
}

func NewAsyncFailoverSMSService(svcs []sms.Service, threshold int32, t time.Duration) *AsyncFailoverSMSService{
	return &AsyncFailoverSMSService{
		svcs: svcs,
		curTime: time.Now().Unix(),
		threshold: threshold,
		timeDuration: t,
	}
	
}


func (t *AsyncFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	prevTime := atomic.LoadInt64(&t.curTime)

	curTime := time.Now()
	// third party server fail
	if cnt>=t.threshold {
		if curTime.Sub(time.Unix(prevTime, 0))<=t.timeDuration{
			newIdx := (idx+1)%int32(len(t.svcs))
			if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx){
				atomic.StoreInt32(&t.cnt, 0)
				atomic.StoreInt64(&t.curTime, curTime.Unix())
			}
			idx = newIdx


		}else{
			atomic.StoreInt32(&t.cnt, 0)
			atomic.StoreInt64(&t.curTime, curTime.Unix())
	
		}
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err{
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		atomic.StoreInt64(&t.curTime, curTime.Unix())
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:

	}
	return err

}
