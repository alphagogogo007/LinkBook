package service

import (
	"context"
	"errors"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
	ErrCodeSendTooMany = cache.ErrCodeSendTooMany
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)

}

// 为什么这里要返回一个domain user？
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)

	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil

}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	// UpdateNicknameAndXXAnd
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *UserService) FindById(ctx context.Context,
	uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) GetUserIdFromSession(ctx *gin.Context) (int64, error) {

	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	if userId == nil {
		return 0, errors.New("wrong session")
	}
	userIdInt64, ok := userId.(int64)
	if !ok {
		return 0, errors.New("wrong type")
	}

	return userIdInt64, nil
}
