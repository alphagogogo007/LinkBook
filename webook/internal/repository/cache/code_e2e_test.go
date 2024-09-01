package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCodeCache_Set_e2e(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name:   "set success",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
				duration, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, duration > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:     context.Background(),
			biz:     "Login",
			phone:   "15212341234",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "send too many",

			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				err := rdb.Set(ctx, key, "123456", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
				duration, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, duration > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:     context.Background(),
			biz:     "Login",
			phone:   "15212341234",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "no expiration",

			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:Login:15212341234"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:     context.Background(),
			biz:     "Login",
			phone:   "15212341234",
			code:    "123456",
			wantErr: errors.New("validation code exist, but has no expiration"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.before(t)
			defer tc.after(t)
			c := NewRedisCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)

		})
	}

}
