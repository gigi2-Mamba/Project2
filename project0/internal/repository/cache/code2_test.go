package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project0/internal/repository/cache/redismocks"
	"testing"
)

func TestRedisCodeCache2_Set(t *testing.T) {
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		name string
		// 模拟真实的依赖项
		mock func(ctrl *gomock.Controller) redis.Cmdable
		// 输入
		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				// 第三方依赖才会这样生成mock吗，并不是每一个依赖都是这么生成的
				res := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc("test", "13168896093")},
					[]any{"123456"}).Return(cmd)
				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "13168896093",
			code:    "123456",
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		//启动testing.T
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			redisCodeCache := NewCodeCache(tc.mock(ctrl))
			err := redisCodeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
