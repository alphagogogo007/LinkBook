package failover

import (
	"context"
	"errors"
	"log"

	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service{
	return &FailoverSMSService{
		svcs: svcs,
	}
}
func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {

	for _,  svc := range f.svcs{
		err := svc.Send(ctx, tplId, args, numbers...)
		if err==nil{
			return nil
		}
		log.Println(err)

	}
	return errors.New("all service failed")
}
