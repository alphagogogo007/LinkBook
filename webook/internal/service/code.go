package service

import (
	"context"
	"fmt"
	"math/rand"

	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
)


type CodeService struct{
	repo *repository.CodeRepository
	sms sms.Service
}

// TODO: no newcodeservice
func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService{
	return &CodeService{
		repo: repo,
		sms:  smsSvc,
	}

}

// 没有new codeService吗？
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error{
	code:= svc.generate()
	err:=svc.repo.Set(ctx, biz, phone, code)
	if err!=nil{
		return err
	}

	const codeTplId =  "12345"
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error){
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err== repository.ErrCodeVerifyTooMany{
		return false, nil
	}
	return ok,err

}

func (svc *CodeService) generate() string{
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}