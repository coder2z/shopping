package cache

import (
	"context"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	rc  *redis.ClusterClient
	one sync.Once
)

func RedisHandle() *redis.ClusterClient {
	one.Do(func() {
		rc = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    xcfg.GetStringSlice("redis.addr"),
			Username: xcfg.GetString("redis.name"),
			Password: xcfg.GetString("redis.pw"),
		})
		ping := rc.Ping(context.Background())
		if ping.Err() != nil {
			panic(ping.Err())
		}
	})
	return rc
}
