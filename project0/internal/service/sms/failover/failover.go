package failover

import (
	"context"
	"errors"
	"project0/internal/service/sms"
	"sync/atomic"
)

type FailOverSMSService struct {
	smss []sms.Service
	idx  uint64
}

func NewFailOverSMSService(smss []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{smss: smss}
}

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//for _,sms := range f.smss {
	//
	//	err := sms.Send(ctx, tplId, args, numbers...)
	//	if err == nil {
	//		// 不return nil怎样，其实就是传nil 可能有微弱的性能优化
	//		return err
	//	}
	//	//return errors.New("发送失败")
	//	log.Println("third manage service",err)
	//}
	//
	//return errors.New("发送失败，所有服务商都尝试过了")
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.smss))

	for i := idx; i < idx+length; i++ {
		svc := f.smss[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)

		switch err {
		case nil:
			//这样可读性强
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}
	}

	return errors.New("发送失败，所有服务商都尝试过了")

}
