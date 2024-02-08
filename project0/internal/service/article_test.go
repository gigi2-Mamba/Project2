package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	reposmocks "project0/.internal/repository/mocks"
	"project0/internal/domain"
	"project0/internal/repository"
	"project0/pkg/loggerDefine"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
     // 做service层的测试就要注入下一层，repository层。 mock repository
	testCases := []struct{
		name string
		mock func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository,repository.ArticleReaderRepository)

		art domain.Article
		wantId int64
		wantErr error

	} {
		{
			name: "新建并发表成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
				authorRepository.EXPECT().Create(gomock.Any(),domain.Article{
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1),nil)
				readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)
	 			readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
					 Id: 1,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				})
				return authorRepository,readerRepository
			},
			art: domain.Article{
				Title: "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: int64(1),
		},
		{
			name: "修改并新发表成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
				authorRepository.EXPECT().Update(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)
				readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				})
				return authorRepository,readerRepository
			},
			art: domain.Article{
				Id: 11,
				Title: "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 11,
		},

		//--实践版本 先不引入重试
		//{
		//	name: "修改并发表失败",
		//	mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
		//		authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
		//		authorRepository.EXPECT().Update(gomock.Any(),domain.Article{
		//			Id: 11,
		//			Title: "我的标题",
		//			Content: "我的内容",
		//			Author: domain.Author{
		//				Id: 123,
		//			},
		//		}).Return(nil)
		//		readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)
		//		readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
		//			Id: 11,
		//			Title: "我的标题",
		//			Content: "我的内容",
		//			Author: domain.Author{
		//				Id: 123,
		//			},
		//		}).Return(errors.New("mock db err"))
		//		return authorRepository,readerRepository
		//	},
		//	art: domain.Article{
		//		Id: 11,
		//		Title: "我的标题",
		//		Content: "我的内容",
		//		Author: domain.Author{
		//			Id: 123,
		//		},
		//	},
		//	wantId: 11,
		//	wantErr:errors.New("mock db err"),
		//},

		{
			name: "修改并发表失败，重试成功",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
				authorRepository.EXPECT().Update(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)
				readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db err"))
				readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return authorRepository,readerRepository
			},
			art: domain.Article{
				Id: 11,
				Title: "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 11,
			wantErr:nil,
		},
		{
			name: "修改并发表失败，重试次数耗尽",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
				authorRepository.EXPECT().Update(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)
				readerRepository.EXPECT().Save(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Times(3).Return(errors.New("mock db err"))

				return authorRepository,readerRepository
			},
			art: domain.Article{
				Id: 11,
				Title: "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId: 11,
			wantErr:errors.New("保存到制作库成功但是线上库失败,重试耗尽"),

		},


		{
			name: "修改并保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				authorRepository := reposmocks.NewMockArticleAuthorRepository(ctrl)
				authorRepository.EXPECT().Update(gomock.Any(),domain.Article{
					Id: 11,
					Title: "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db err"))
				readerRepository := reposmocks.NewMockArticleReaderRepository(ctrl)

				return authorRepository,readerRepository
			},
			art: domain.Article{
				Id: 11,
				Title: "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr:errors.New("mock db err"),
		},
	}


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authorRepository, readerRepository := tc.mock(ctrl)
			svc := NewArticleServiceV1(authorRepository, readerRepository ,loggerDefine.NewNopLogger())

            id,err := svc.PublishV1(context.Background(),tc.art)
			assert.Equal(t, tc.wantErr,err)
			assert.Equal(t, tc.wantId,id)

		})
	}
}
