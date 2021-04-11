package utils

import (
	"context"
	"shopping/cache"
)

const (
	periodScript = `
local num = redis.call("GET", KEYS[1])
if num == 0 then
    return false
end

local current = redis.call("INCRBY", KEYS[1], -1)
if current < 0 then
    return false
else
    return true
end`
)

func Limit(ctx context.Context, key string) bool {
	eval := cache.RedisHandle().Eval(ctx, periodScript, []string{key}, []string{})
	ok, err := eval.Bool()
	if err != nil {
		Log.Error("limit error", err.Error())
	}
	return ok
}

func AddStock(ctx context.Context, key string, num int64) error {
	return cache.RedisHandle().Set(ctx, key, num, 0).Err()
}

func StockAddOne(ctx context.Context, key string) error {
	return cache.RedisHandle().Incr(ctx, key).Err()
}
