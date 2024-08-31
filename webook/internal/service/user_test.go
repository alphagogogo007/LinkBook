package service

import (
	"context"
	"errors"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	assert.NoError(t, err)
}

func TestRegularUserService_Login(t *testing.T) {

	userEmail1 := "123@qq.com"
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		args     args
		wantUser domain.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name: "login successful",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), userEmail1).Return(
					domain.User{
						Email:    userEmail1,
						Password: "$2a$10$pzhe5saJTm7yQIU52dM5fu1ZzSjlUwI/RocB79zmqK1LytKx9IK8K",
						Phone: "15212341234",
					}, nil)
				return repo

			},
			args: args{
				ctx:      context.Background(),
				email:    userEmail1,
				password: "12345678",
			},
			wantUser: domain.User{
				Email:   	userEmail1,
				Password: "$2a$10$pzhe5saJTm7yQIU52dM5fu1ZzSjlUwI/RocB79zmqK1LytKx9IK8K",
				Phone: "15212341234",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), userEmail1).Return(
					domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			args: args{
				ctx:      context.Background(),
				email:    userEmail1,
				password: "12345678",
			},
			wantUser: domain.User{
			},
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), userEmail1).Return(
					domain.User{}, errors.New("db error"))
				return repo
			},
			args: args{
				ctx:      context.Background(),
				email:    userEmail1,
				password: "12345678",
			},
			wantUser: domain.User{
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "wrong password",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), userEmail1).Return(
					domain.User{
						Email:    userEmail1,
						Password: "$2a$10$pzhe5saJTm7yQIU52dM",
						Phone: "15212341234",
					}, nil)
				return repo

			},
			args: args{
				ctx:      context.Background(),
				email:    userEmail1,
				password: "12345678",
			},
			wantUser: domain.User{},
			wantErr: ErrInvalidUserOrPassword,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tt.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tt.args.ctx, tt.args.email, tt.args.password)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUser, user)

		})
	}
}
