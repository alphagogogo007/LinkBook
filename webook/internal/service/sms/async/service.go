package async

import (
	"context"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"go.uber.org/zap"
	"math/rand"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
)

type Service struct {
	svc  sms.Service
	repo repository.AsyncSmsRepository
	l    *zap.Logger //what is this logger?
	cnt  int32
	cntThreshold int32
	timeoutThreshold time.Duration
	limiter limiter.Limiter
	key     string
}

func NewService(svc sms.Service, repo repository.AsyncSmsRepository, l *zap.Logger, 
	cnt int32, cntThreshold int32, timeoutThreshold time.Duration, limiter limiter.Limiter) *Service {
	res := &Service{
		svc:  svc,
		repo: repo,
		l:l,
		cnt: cnt,
		cntThreshold: cntThreshold,
		timeoutThreshold: timeoutThreshold,
		limiter: limiter,
		key:   "async-limiter",
	}
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

func (s *Service) StartAsyncCycle() {

	//avoid panic
	defer func() {
        if r := recover(); r != nil {
            s.l.Error("panic in StartAsyncCycle", zap.Any("error", r))
            go s.StartAsyncCycle() 
        }
    }()

	time.Sleep(time.Second * 3)
	for {
		s.AsyncSend()
	}
}

func (s *Service) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptWaitingSMS(ctx) //这里为什么就抢占了？

	cancel() //这里为什么有个cancel？因为要提前结束
	switch err {
	case nil:
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		if err != nil {
			s.l.Error("tried to send, but failed", zap.Error(err), zap.Int64("Id", as.Id)) //这个logger怎么使用？
		}
		res := err == nil
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if err != nil {

			s.l.Error("mark database error",
				zap.Error(err),
				zap.Bool("res", res),
				zap.Int64("Id", as.Id))
		}

	case repository.ErrWaitingSMSNotFound:
		time.Sleep(time.Second)
	default:
		//database error
		s.l.Error("async send failure", zap.Error(err)) //这里要补充pkg包
		time.Sleep(time.Second)
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.needAsync() {
		err := s.repo.Add(ctx, domain.AsyncSms{
			TplId:    tplId,
			Args:     args,
			Numbers:  numbers,
			RetryMax: 3,
		})
		return err
	}
	start := time.Now()
	err := s.svc.Send(ctx, tplId, args, numbers...) //这里只有一个svc
	duration := time.Since(start)
	if duration>=s.timeoutThreshold{
		s.cnt += 1
	}else if err==nil{
		s.cnt = 0
	}
	return err
}

// write my logic
func (s *Service) needAsync() bool {
	
	if s.cnt>=s.cntThreshold{
		// 连续10次超时响应后，剩下的request随机同步和异步发送。
		//比如10%同步发送，如果成功，cnt重置为0.不成功，cnt+1
		// 90%异步发送的不参加计数。
		rnd := rand.Intn(100)
		return !(rnd<10)
	}
	limited, err := s.limiter.Limit(context.Background(), s.key)
	if err != nil {
		s.l.Error("async limiter error")
		return false
	}
	if limited {
	// ratelimiter
		s.l.Warn("trigger async limiter")
		return true
	}
	return false

}
