package failover

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	smsmocks "project0/.internal/service/sms/mocks"
	"project0/internal/service/sms"
	"testing"
)

func TestTimeFailoverSMSService_Send(t1 *testing.T) {

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) []sms.Service

		threshold int32
		idx       int32
		cnt       int32

		//wantIdx int32
		//wantCnt int32
		wantErr error
	}{
		{
			name: "没有切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0}
			},
			idx:       0,
			cnt:       12,
			threshold: 15,

			//wantIdx:   0,
			// 成功了，重置超时计数
			//wantCnt: 0,
			wantErr: nil,
		},
		{
			name: "触发切换成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				svc0 := smsmocks.NewMockService(ctrl)

				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,

			wantErr: nil,
			//wantCnt: 0,
			//wantIdx: 1,
		},
		{
			name: "触发切换成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			wantErr:   nil,
			//wantCnt: 0,
			//wantIdx: 1,
		},
	}

	for _, tc := range testCases {
		t1.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t1)
			defer ctrl.Finish()
			//timeFailoverSMSService := tc.mock(ctrl)
			timeouversvc := NewTimeFailoverSMSService(tc.mock(ctrl), tc.threshold)
			timeouversvc.idx = tc.idx
			timeouversvc.cnt = tc.cnt
			err := timeouversvc.Send(context.Background(), "abcde", []string{"dsfds"}, "123242cdfg")
			assert.Equal(t, err, tc.wantErr)
			//assert.Equal(t, timeouversvc.cnt,tc.cnt)
			//assert.Equal(t, timeouversvc.idx,tc.idx)
		})
	}
}
