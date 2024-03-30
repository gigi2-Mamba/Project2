package ratelimit

import (
	"context"
	"errors"
	"log"
	"project0/internal/service/sms"
	"project0/pkg/limiter"
)

var errLimited = errors.New("sms触发限流")

type RateLimitSMSService struct {
	// 被装饰的
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

// 基于接口实现的 可以找到顶层的接口实现然后mock
func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		log.Println("redis 限流sms 有误")
		return err
	}
	if limited {
		return errLimited
	}

	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
		key:     "sms-limiter",
	}
}
