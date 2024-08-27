package ioc

import (
	"context"
	"log"

	"gitee.com/geekbang/basic-go/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	// 检查 Redis 连接是否成功
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return redisClient
}
