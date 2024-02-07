package ratelimit

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

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		//  一个没有要测的逻辑和输入没关系
		//ctx context.Context
		//phones []string
		//tpId string
		//args []string

		//
		wantErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter0 := limitermocks.NewMockLimiter(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				limiter0.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				return svc, limiter0

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc, limiter0 := tc.mock(ctrl)
			rservice := NewRateLimitSMSService(svc, limiter0)
			err := rservice.Send(context.Background(), "sss", []string{"sss"}, "1232d")
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
