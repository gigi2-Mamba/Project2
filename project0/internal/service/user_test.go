package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	reposmocks "project0/.internal/repository/mocks"
	"project0/internal/domain"
	"project0/internal/repository"

	"testing"
)

func Test_userService_Signup(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期输入
		context  context.Context
		email    string
		password string

		//预期输出
		wantUser domain.User
		wantErr  error
	}{
		{name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repos := reposmocks.NewMockUserRepository(ctrl)
				// mock出来的 userRepository 指明特定方法返回的期待值
				repos.EXPECT().FindByEmail(gomock.Any(), "8@qq.com").
					Return(domain.User{
						Email:    "8@qq.com",
						Password: "$2a$10$sO8Lw6ukyzdC6XZsVSwfh.JjpvBPlUAqlySeVbyaxIM7FCXGls8TO",
						Phone:    "13168896093",
					}, nil)
				return repos
			},
			email:    "8@qq.com",
			password: "123456@q",

			wantUser: domain.User{
				Email:    "8@qq.com",
				Password: "$2a$10$sO8Lw6ukyzdC6XZsVSwfh.JjpvBPlUAqlySeVbyaxIM7FCXGls8TO",
				Phone:    "13168896093",
			},
			// wanterr 不写默认就是nil
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRepository := tc.mock(ctrl)
			// 这里测试要使用具体实现
			svc := NewUserService(userRepository)
			user, err := svc.Login(tc.context, tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)

		})
	}

}
