package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	daomocks "project0/.internal/repository/dao/mocks"
	"project0/internal/domain"
	"project0/internal/repository/cache"
	"project0/internal/repository/cache/mocks"
	"project0/internal/repository/dao"
	"testing"
	"time"
)

func TestCacheUserRepository_Profile(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)

		//预期输入
		ctx context.Context
		uid int64

		//输出
		wantUserProfile domain.UserProfile
		wantErr         error
	}{
		{
			name: "缓存未命中,读取成功",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(3)
				cachemock := cachemocks.NewMockUserCache(ctrl)
				daomock := daomocks.NewMockUserDao(ctrl)
				cachemock.EXPECT().Get(gomock.Any(), uid).Return(domain.UserProfile{}, cache.ErrKeyNotExist)
				daomock.EXPECT().Profile(gomock.Any(), uid).Return(dao.UserProfile{
					Id:           3,
					Gender:       "F",
					NickName:     "waste time",
					Introduction: "introduce",
					BirthDate:    time.UnixMilli(101).UnixMilli(),
				}, nil)
				cachemock.EXPECT().Set(gomock.Any(), domain.UserProfile{
					Id:           3,
					Gender:       "F",
					NickName:     "waste time",
					Introduction: "introduce",
					BirthDate:    time.UnixMilli(101)}).Return(nil)
				return cachemock, daomock
			},
			ctx: context.Background(),
			uid: 3,

			wantUserProfile: domain.UserProfile{
				Id:           3,
				Gender:       "F",
				NickName:     "waste time",
				Introduction: "introduce",
				BirthDate:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create controller first
			controller := gomock.NewController(t)
			// mock了依赖项
			userCache, userDao := tc.mock(controller)
			svc := NewCacheUserRepository(userDao, userCache)
			userProfile, err := svc.Profile(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantUserProfile, userProfile)
			assert.Equal(t, tc.wantErr, err)
			return

		})
	}
}
