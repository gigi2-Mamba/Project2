package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

// 得在var的括号里面定义
var (
	//go:embed redis_slide_window.lua
	luaScript string
)

type RedisSlideWindowLimiter struct {
	cmd redis.Cmdable

	interval time.Duration
	rate     int
}

func (r *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	// 配合lua的go-redis 命令调用
	return r.cmd.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlideWindowLimiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}
