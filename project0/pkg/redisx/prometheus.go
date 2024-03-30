package redisx

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

/*
Created by payor-programmer on 2024/2/18.
*/

/*
利用redis的原生hook实现埋点
*/

type Prometheus struct {
	vector *prometheus.SummaryVec
}

func NewPrometheus(opt prometheus.SummaryOpts) *Prometheus {
	return &Prometheus{
		vector: prometheus.NewSummaryVec(opt, []string{
			"cmd", "key_exist",
		}),
	}
}

func (p Prometheus) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (p Prometheus) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		// 其实就是为了统计使用的命令，是否存在，响应时间
		start := time.Now()
		var err error
		defer func() {
			duration := time.Since(start).Milliseconds()
			keyExists := err == redis.Nil
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExists)).Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		return err

	}
}

func (p Prometheus) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
