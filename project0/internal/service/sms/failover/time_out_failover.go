package failover

import (
	"context"
	"project0/internal/service/sms"
	"sync/atomic"
)

type TimeFailoverSMSService struct {
	svcs []sms.Service
	// 当前正在使用节点
	idx int32
	// 连续几个超时
	cnt int32
	// 切换的阈值
	threshold int32
}

func NewTimeFailoverSMSService(svcs []sms.Service, threshold int32) *TimeFailoverSMSService {
	return &TimeFailoverSMSService{svcs: svcs, threshold: threshold}
}

// 判断服务商是否真的可用，采用的第三方服务切换
// 严格超时需要使用锁机制
func (t *TimeFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	// 如果连续超时超过阈值。 天然不需要修改的数据，没有并发问题
	// 当时只写了大于，会报错  unexpected call
	if cnt >= t.threshold {
		//触发切换，切换新的服务商+1
		neweIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, neweIdx) {
			// 切换成功，需要把超时计数重置为0
			atomic.StoreInt32(&t.cnt,0)

		}
		// 为了给后面调用写的
		idx = neweIdx

	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)

	switch err {
	case nil:
		//没有任何错误，重置计数器
		atomic.CompareAndSwapInt32(&t.cnt, cnt, 0)
		// 这里不走也可以啊，可以依赖后面的
		return nil
	case context.DeadlineExceeded:
		// 仅仅是超时，没有其他的了
		atomic.AddInt32(&t.cnt, 1)
		//return err
	default:
		//不是服务超时错误，但是错误，怎么搞
		// 当做超时错误还是收走原子直接切换
		// EOF错误直接切换
	}
	return err
}
