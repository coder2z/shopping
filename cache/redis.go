package cache

import (
	"context"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	rc  *redis.Client
	one sync.Once
)

func RedisHandle() *redis.Client {
	one.Do(func() {
		rc = redis.NewClient(&redis.Options{
			Addr:     xcfg.GetString("redis.addr"),
			Username: xcfg.GetString("redis.name"),
			Password: xcfg.GetString("redis.pw"),
			DB:       xcfg.GetInt("redis.db"),
		})
		ping := rc.Ping(context.Background())
		if ping.Err() != nil {
			panic(ping.Err())
		}
	})
	return rc
}
