package cache

import (
	"context"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode        string
	ErrCodeSendTooMany   = errors.New("send too many")
	ErrCodeVerifyTooMany = errors.New("verify too many")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
	key(biz, phone string) string
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

type BigCacheCodeCache struct {
	cache *bigcache.BigCache
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case -2:
		return errors.New("validation code exist, but has no expiration")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil

	}

}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil

	}

}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func NewBigCacheCodeCache(cache *bigcache.BigCache) CodeCache {
	return &BigCacheCodeCache{
		cache: cache,
	}
}

func (c *BigCacheCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := c.key(biz, phone)
	val, err := c.cache.Get(key)
	if err == nil {
		// 获取存储时间戳并解析
		parts := string(val)
		storedTime := time.Unix(0, int64(binary.LittleEndian.Uint64([]byte(parts[:8]))))
		if time.Since(storedTime) < 9*time.Minute {
			return ErrCodeSendTooMany // 如果时间小于9分钟，返回发送过多错误
		}
	}

	// 记录当前时间戳，并将其与验证码一起存入缓存
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(time.Now().UnixNano()))
	data := append(timestamp, []byte(code)...)

	return c.cache.Set(key, data)
}


func (c *BigCacheCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	key := c.key(biz, phone)
	val, err := c.cache.Get(key)
	if err != nil {
		return false, err
	}

	if string(val[8:]) != code {
		return false, nil // 如果验证码不匹配，返回验证过多错误
	}

	// 验证通过后可以选择是否删除缓存项（验证码一次性使用）
	_ = c.cache.Delete(key)

	return true, nil
}

func (c *BigCacheCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
