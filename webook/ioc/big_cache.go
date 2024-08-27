package ioc

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

func InitBigCache(ctx context.Context) *bigcache.BigCache {

	config := bigcache.DefaultConfig(10 * time.Minute) // 过期时间配置为10分钟
	config.Shards = 1024                               // 分片数量
	config.MaxEntriesInWindow = 1000 * 10 * 60         // 每个窗口的最大条目数
	config.MaxEntrySize = 500                          // 每个缓存项的最大大小（字节）
	config.Verbose = true                              // 输出详细日志

	cache, err := bigcache.New(ctx, config)
	if err != nil {
		panic(err)
	}
	return cache
}
