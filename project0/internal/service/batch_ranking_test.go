package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project0/internal/domain"
	svcmocks "project0/internal/service/mocks"
	"testing"
	"time"
)

/*
Created by society-programmer on 2024/2/20.
*/

func TestRankingService(t *testing.T)  {
	const batchSize = 2
	now := time.Now()
    testCases :=[]struct{
		name string
		mock func(ctrl *gomock.Controller) (InteractiveService,ArticleService)

		wantErr error
		wantArts []domain.Article
	}{
		{
			name: "成功获取",
			mock: func(ctrl *gomock.Controller) (InteractiveService, ArticleService) {
				articleSvc := svcmocks.NewMockArticleService(ctrl)
				intersvc :=svcmocks.NewMockInteractiveService(ctrl)
				articleSvc.EXPECT().ListPub(gomock.Any(),gomock.Any(),0,2).
					Return([]domain.Article{
						{Id: 1, Utime: now},
						{Id: 2, Utime: now},
				},nil)
				//第二批
				articleSvc.EXPECT().ListPub(gomock.Any(),gomock.Any(),2,2).
					Return([]domain.Article{
						{Id: 3,Utime: now},
						{Id: 4,Utime: now},
					},nil)
				//第三批没有数据
				articleSvc.EXPECT().ListPub(gomock.Any(),gomock.Any(),4,2).
					Return([]domain.Article{},nil)

				intersvc.EXPECT().GetByIds(gomock.Any(),"article",[]int64{1,2}).
					Return(map[int64]domain.Interactive{
						1:{LikeCnt: 1},
						2: {LikeCnt: 2},
				},nil)

				// 第二批
				intersvc.EXPECT().GetByIds(gomock.Any(),"article",[]int64{3,4}).
					Return(map[int64]domain.Interactive{
						3:{LikeCnt: 3},
						4: {LikeCnt: 4},
					},nil)

				// 第三批的点赞数据
				intersvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).
					Return(map[int64]domain.Interactive{}, nil)

				return intersvc,articleSvc
			},
            wantErr: nil,
			wantArts: []domain.Article{
				{Id: 4, Utime: now},
				{Id: 3, Utime: now},
				{Id: 2, Utime: now},
	},
	},}



	for _,tc := range testCases {
		//启动具体测试的函数
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			inter,article := tc.mock(ctrl)
			//batchSvc := NewBatchRankingService(inter,article)
			batchSvc := &BatchRankingService{
				interSvc:  inter,
				artSvc:   article,
				batchSize: batchSize,
				n:         3,
				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
				return float64(likeCnt)
			},
				}
			res, err := batchSvc.topN(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tc.wantArts,res)

		})
	}

}
