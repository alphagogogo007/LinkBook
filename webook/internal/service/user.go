package service

import (
	"context"
	"errors"
	"fmt"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func (svc *UserService) Edit(ctx context.Context, userProfile domain.UserProfile )  error{

	existUser, err := svc.repo.FindProfileById(ctx, userProfile.UserId)
	fmt.Println("find error", err)
	if err==repository.ErrUserNotFound{
		err = svc.repo.CreateProfile(ctx, userProfile)
		return err
		
	}else if err!=nil{
		return  err
	}
	fmt.Println(existUser, "existing")

	existUser.NickName = userProfile.NickName
	existUser.Birthday = userProfile.Birthday
	existUser.AboutMe = userProfile.AboutMe
	err = svc.repo.OverwriteProfle(ctx, existUser)
	return err


}


func (svc *UserService) GetProfile(ctx context.Context, userId int64)  (domain.FrontProfile, error){
	existUser, err := svc.repo.FindProfileById(ctx, userId)
	if err!=nil{
		return domain.FrontProfile{}, err
	}
	
	return domain.FrontProfile{
		NickName: existUser.NickName,
		Birthday: existUser.Birthday,
		AboutMe: existUser.AboutMe,
	},nil
}

func (svc *UserService) GetUserIdFromSession(ctx *gin.Context) (int64, error){

	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	if userId==nil{
		return 0, errors.New("Didn't get user Id")
	}
	userIdInt64, ok := userId.(int64)
	if !ok {
		return 0, errors.New("user Id wrong type")
	}

	return userIdInt64,nil
}