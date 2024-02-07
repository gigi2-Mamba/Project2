package failover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	smsmocks "project0/.internal/service/sms/mocks"
	"project0/internal/service/sms"
	"testing"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc}
			},
		},
		{name: "第二次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc, svc0}
			},
		},
		{name: "第二次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return []sms.Service{svc, svc0}
			},
			wantErr: errors.New("发送失败，所有服务商都尝试过了"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctrl.Finish()
			//services := smsmocks.NewMockService(ctrl)
			failoverService := NewFailOverSMSService(tc.mock(ctrl))
			err := failoverService.Send(context.Background(), "123", []string{"1232"}, "12345")
			assert.Equal(t, err, tc.wantErr)

		})
	}
}
