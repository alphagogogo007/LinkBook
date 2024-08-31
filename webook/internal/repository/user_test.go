package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	cachemocks "gitee.com/geekbang/basic-go/webook/internal/repository/cache/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {

	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)
		ctx      context.Context
		uid      int64
		wantUser domain.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name: "no cache",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(
					dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "12345678",
						Birthday: 123,
						Nickname: "",
						AboutMe:  "",
						Phone: sql.NullString{
							String: "123",
							Valid:  true,
						},
						CreateAt: 100,
						UpdateAt: 101,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "123",
					Password: "12345678",
					Nickname: "",
					Birthday: time.UnixMilli(123),
					AboutMe:  "",
				}).Return(nil)

				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Phone:    "123",
				Password: "12345678",
				Nickname: "",
				Birthday: time.UnixMilli(123),
				AboutMe:  "",
			},
			wantErr: nil,
		},
		{
			name: "find cache",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "123",
					Password: "12345678",
					Nickname: "",
					Birthday: time.UnixMilli(123),
					AboutMe:  "",
				}, nil)

				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Phone:    "123",
				Password: "12345678",
				Nickname: "",
				Birthday: time.UnixMilli(123),
				AboutMe:  "",
			},
			wantErr: nil,
		},
		{
			name: "no user",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(
					dao.User{}, dao.ErrRecordNotFound)
				

				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
			},
			wantErr:ErrUserNotFound,
		},
		{
			name: "write cache error",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(
					dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "12345678",
						Birthday: 123,
						Nickname: "",
						AboutMe:  "",
						Phone: sql.NullString{
							String: "123",
							Valid:  true,
						},
						CreateAt: 100,
						UpdateAt: 101,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "123",
					Password: "12345678",
					Nickname: "",
					Birthday: time.UnixMilli(123),
					AboutMe:  "",
				}).Return(errors.New("redis error"))

				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Phone:    "123",
				Password: "12345678",
				Nickname: "",
				Birthday: time.UnixMilli(123),
				AboutMe:  "",
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc, ud := tt.mock(ctrl)
			repo := NewUserRepository(ud, uc)

			user, err := repo.FindById(tt.ctx, tt.uid)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUser, user)

		})
	}
}
