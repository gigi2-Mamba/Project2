package failover

import (
	"context"
	"log"
	"math"
	"project0/internal/service/sms"
	"project0/pkg/limiter"
	"sync/atomic"
	"time"
)

/*
判断服务不可用
本次超时机制超过限定的响应时间或者这次响应时间到达上次的2.5倍
*/
type ResponseTimeFailover struct {
	smss      []sms.Service
	respTime  int64
	limiter   limiter.Limiter
	limited   bool
	// 最大接口响应时间
	threshold int64
	// 当前节点服务商
	idx int32
	//限流的key
	key string
	// 允许接口响应时间的倍数
	diff int64
}

func newReq(ctx context.Context, tplId string, args []string, idx int32, numbers ...string) *AsyncSendCodeReq {

	return &AsyncSendCodeReq{
		Ctx:     ctx,
		TplId:   tplId,
		Args:    args,
		Idx:     idx,
		Numbers: numbers,
	}
}
func (r *ResponseTimeFailover) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	start := time.Now()
	idx := atomic.LoadInt32(&r.idx)
	//log.Println("current Call Send,r.idx is ",idx)
	lastSpent := atomic.LoadInt64(&r.respTime)
	limited0, err := r.limiter.Limit(ctx, r.key)
	if err != nil || limited0 {
		r.limited = true
	}
	// 达到限流限制，同步转异步
	if r.limited {
		log.Println("发送验证码频繁触发限流,当前下表为： ",idx)
		req := newReq(ctx, tplId, args, idx+1, numbers...)
		ReqChan <- req
		return nil

	}


	err = r.smss[idx].Send(ctx, tplId, args, numbers...)

	spendTime := int64(math.Ceil(time.Now().Sub(start).Seconds()))

	log.Printf("spendtime is %v, lastSpent is %v",spendTime,lastSpent)
	atomic.StoreInt64(&r.respTime,spendTime)
	// 服务商崩溃 // 服务商崩溃响应时间超过两秒
	if spendTime >= r.threshold || spendTime / lastSpent == r.diff {
		log.Println("服务商崩溃")
		req := newReq(ctx, tplId, args, idx+1, numbers...)
		ReqChan <- req
	}
	return err
}

func NewResponseTimeFailover(smss []sms.Service, limiter limiter.Limiter, threshold, diff int64, key string) *ResponseTimeFailover {
	return &ResponseTimeFailover{
		smss:      smss,
		limiter:   limiter,
		threshold: threshold,
		key:       key,
		diff: diff,
		respTime: diff,
	}
}
