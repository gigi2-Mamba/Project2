package ioc

import (
	"github.com/redis/go-redis/v9"
	"project0/pkg/limiter"
	"time"
)

// for responseFailover temporary use
func InitRedisLimiter(cmd redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlideWindowLimiter(cmd, time.Second*5, 2)
}
