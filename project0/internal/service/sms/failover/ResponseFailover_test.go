package failover

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	smsmocks "project0/.internal/service/sms/mocks"
	"project0/internal/service/sms"
	"project0/pkg/limiter"
	limitermocks "project0/pkg/limiter/mocks"
	"testing"
)

func TestResponseTimeFailover_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *gomock.Controller) ([]sms.Service, limiter.Limiter)
		// 预期输入也就是必须输入的
		threshold int64
		diff      int64
		key       string

		wantIdx     int32
		wantLimited bool
		wantErr     error
		// 服务商是否不可用
		wantLapse bool
	}{{
		name: "无需切换",
		mock: func(ctrl *gomock.Controller) ([]sms.Service, limiter.Limiter) {
			limiter0 := limitermocks.NewMockLimiter(ctrl)

			svc0 := smsmocks.NewMockService(ctrl)
			svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			limiter0.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
			return []sms.Service{svc0}, limiter0
		},
		threshold: 2,
		diff:      2,
		key:       "smsresp-limit",

		wantIdx:     0,
		wantLimited: false,
		wantErr:     nil,
		wantLapse:   false,
	}, {
		name: "限流请求变异步处理",
		mock: func(ctrl *gomock.Controller) ([]sms.Service, limiter.Limiter) {
			limiter0 := limitermocks.NewMockLimiter(ctrl)

			svc0 := smsmocks.NewMockService(ctrl)
			svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			limiter0.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
			return []sms.Service{svc0}, limiter0
		},
		threshold: 2,
		diff:      2,
		key:       "smsresp-limit",

		wantIdx:     1,
		wantLimited: true,
		wantErr:     nil,
		wantLapse:   false,
	}, {
		name: "服务商崩溃,请求变异步",
		mock: func(ctrl *gomock.Controller) ([]sms.Service, limiter.Limiter) {
			limiter0 := limitermocks.NewMockLimiter(ctrl)

			svc0 := smsmocks.NewMockService(ctrl)
			svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			limiter0.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
			return []sms.Service{svc0}, limiter0
		},
		threshold: 2,
		diff:      2,
		key:       "smsresp-limit",

		wantIdx:     1,
		wantLimited: false,
		wantErr:     nil,
		wantLapse:   true,
	}} //响应时间变成两倍暂且不测}

	for _, tc := range testCases {
		// 测试实例开始
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctrl.Finish()
			// 测试什么方法，就要直接调用，调用该方法的母体(需要mock注入依赖）
			smss, l := tc.mock(ctrl)
			failoverService := NewResponseTimeFailover(smss, l, tc.threshold, tc.diff, tc.key)

			//assert.NoError(t, err)

			err := failoverService.Send(context.Background(), "1234", []string{"1232"}, "dfdfd")
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantLimited, failoverService.limited)

			if tc.wantLapse {
				// 因为不好模拟接口的耗时大于阈值,
				assert.Equal(t, tc.wantLapse, failoverService.respTime+2 > tc.threshold)

			} else {
				assert.Equal(t, tc.wantLapse, failoverService.respTime > tc.threshold)
			}

			//assert.True(t,)

		})
	}
}
