package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCodeCache_Set(t *testing.T) {

	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}

	type args struct {
		ctx   context.Context
		biz   string
		phone string
		code  string
	}
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "set success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdResult := redis.NewCmdResult(int64(0), nil)
				res := redismocks.NewMockCmdable(ctrl)
				// ...interface{}这里对应的要用[]any{}
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "12312341234")}, []any{"123456"} ).Return(cmdResult)
				return res
			},
			args: args{
				ctx:   context.Background(),
				biz:   "test",
				phone: "12312341234",
				code:  "123456",
			},
			wantErr: nil,
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdResult := redis.NewCmdResult(int64(0), errors.New("redis error"))
				res := redismocks.NewMockCmdable(ctrl)
				// ...interface{}这里对应的要用[]any{}
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "12312341234")}, []any{"123456"} ).Return(cmdResult)
				return res
			},
			args: args{
				ctx:   context.Background(),
				biz:   "test",
				phone: "12312341234",
				code:  "123456",
			},
			wantErr: errors.New("redis error"),
		},
		{
			name: "no expiration",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdResult := redis.NewCmdResult(int64(-2), nil)
				res := redismocks.NewMockCmdable(ctrl)
				// ...interface{}这里对应的要用[]any{}
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "12312341234")}, []any{"123456"} ).Return(cmdResult)
				return res
			},
			args: args{
				ctx:   context.Background(),
				biz:   "test",
				phone: "12312341234",
				code:  "123456",
			},
			wantErr:  errors.New("validation code exist, but has no expiration"),
		},
		{
			name: "send too many",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdResult := redis.NewCmdResult(int64(-1), nil)
				res := redismocks.NewMockCmdable(ctrl)
				// ...interface{}这里对应的要用[]any{}
				res.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "12312341234")}, []any{"123456"} ).Return(cmdResult)
				return res
			},
			args: args{
				ctx:   context.Background(),
				biz:   "test",
				phone: "12312341234",
				code:  "123456",
			},
			wantErr:   ErrCodeSendTooMany,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tt.mock(ctrl))
			err := c.Set(tt.args.ctx, tt.args.biz, tt.args.phone, tt.args.code)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}
