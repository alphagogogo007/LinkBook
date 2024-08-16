package service

import (
	"context"
	"errors"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService struct{
	repo *repository.UserRepository
}

func NewUserService(repo  *repository.UserRepository) *UserService{
	return &UserService{
		repo: repo,
	}
}


func (svc *UserService) SignUp(ctx context.Context, user domain.User ) error{
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err!= nil{
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)

} 

// 为什么这里要返回一个domain user？
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)

	if err==repository.ErrUserNotFound{
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err!=nil{
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password) )
	if err!=nil{
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil

}