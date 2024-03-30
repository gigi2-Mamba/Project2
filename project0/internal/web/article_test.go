package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	service2 "project0/interactive/service"
	"project0/internal/domain"
	"project0/internal/service"
	svcmocks "project0/internal/service/mocks"
	"project0/internal/web/ijwt"
	"project0/pkg/loggerDefine"
	"testing"
)

// 单元测试做handler
//func TestArticleHandler2_Publish(t *testing.T) {
//	testCases := []struct {
//		name string
//		mock func(ctrl *gomock.Controller) service.ArticleService
//
//		reqBody  string
//		wantCode int
//		wantRes  Result
//	}{
//		{
//			name: "修改并且发表成功",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Id:      123,
//					Title:   "新的标题1",
//					Content: "新的内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(123), nil)
//				return svc
//			},
//			reqBody: `
//{
//"id": 123,
// "title": "新的标题1",
// "content": "新的内容"
//}
//`,
//			wantCode: 200,
//			wantRes: Result{
//				// 原本是 int64的，但是因为 Data 是any，所以在反序列化的时候，
//				// 用的 float64
//				Data: float64(123),
//			},
//		},
//		{
//			name: "新建并发表成功x",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Title:   "够标题",
//					Content: "mock",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(1), nil)
//				return svc
//			},
//			reqBody:  `{ "title":"够标题","content":"mock"}`,
//			wantCode: http.StatusOK,
//			wantRes:  Result{Data: float64(1)},
//		},
//		{
//			name: "修改并且发表成功2",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Id: 123, //有没有都一样？
//					// 一旦title和content与数据库不一致就出错
//					//Title: "我的标题",
//					//---
//					//Title: "我的标题",
//					//Content: "我的内容2",
//					Title:   "新的标题1",
//					Content: "新的内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(123), nil)
//				return svc
//			},
//			reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
//			wantCode: http.StatusOK,
//			//在json格式中any类型默认使用int64
//			wantRes: Result{Data: float64(123)},
//		},
//		{
//			name: "已有帖子并发表失败test",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Id:      123,
//					Title:   "新的标题1",
//					Content: "新的内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(123), nil)
//				return svc
//			},
//			//reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
//			reqBody: `
//{
//"id": 123,
// "title": "新的标题1",
// "content": "新的内容"
//}
//`,
//			wantCode: 200,
//			wantRes:  Result{Data: float64(123)},
//		},
//		{
//			name: "已有帖子并发表失败",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					//Id:      123,
//					Title:   "新的标题1",
//					Content: "新的内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(0), errors.New("系统错误"))
//				return svc
//			},
//			//reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
//			reqBody: `
//{
// "title": "新的标题1",
// "content": "新的内容"
//}
//`,
//			wantCode: 200,
//			wantRes: Result{
//				Code: 5,
//				Msg:  "系统错误"},
//		}}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			// 构造 handler
//			svc := tc.mock(ctrl)
//			hdl := NewArticleHandler(svc, loggerDefine.NewNopLogger())
//
//			// 准备服务器，注册路由
//			server := gin.Default()
//			server.Use(func(ctx *gin.Context) {
//				ctx.Set("user", ijwt.UserClaims{
//					Uid: 123,
//				})
//			})
//			hdl.RegisterRoutes(server)
//
//			// 准备Req和记录的 recorder
//			req, err := http.NewRequest(http.MethodPost,
//				"/articles/publish", bytes.NewReader([]byte(tc.reqBody)))
//			assert.NoError(t, err)
//			req.Header.Set("Content-Type", "application/json")
//			recorder := httptest.NewRecorder()
//
//			// 执行
//			server.ServeHTTP(recorder, req)
//			// 断言结果
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			if recorder.Code != http.StatusOK {
//				return
//			}
//			var res Result
//			err = json.NewDecoder(recorder.Body).Decode(&res)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.wantRes, res)
//
//		})
//	}
//}


func TestArticleHandler2_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService)

		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "修改并且发表成功",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      123,
					Title:   "新的标题1",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(123), nil)

				return svc,intrSvc
			},
			reqBody: `
{
"id": 123,
 "title": "新的标题1",
 "content": "新的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				// 原本是 int64的，但是因为 Data 是any，所以在反序列化的时候，
				// 用的 float64
				Data: float64(123),
			},
		},
		{
			name: "新建并发表成功x",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "够标题",
					Content: "mock",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc,intrSvc
			},
			reqBody:  `{ "title":"够标题","content":"mock"}`,
			wantCode: http.StatusOK,
			wantRes:  Result{Data: float64(1)},
		},
		{
			name: "修改并且发表成功2",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService)  {
				svc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id: 123, //有没有都一样？
					// 一旦title和content与数据库不一致就出错
					//Title: "我的标题",
					//---
					//Title: "我的标题",
					//Content: "我的内容2",
					Title:   "新的标题1",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(123), nil)
				return svc,intrSvc
			},
			reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
			wantCode: http.StatusOK,
			//在json格式中any类型默认使用int64
			wantRes: Result{Data: float64(123)},
		},
		{
			name: "已有帖子并发表失败test",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService)  {
				svc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      123,
					Title:   "新的标题1",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(123), nil)
				return svc,intrSvc
			},
			//reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
			reqBody: `
{
"id": 123,
 "title": "新的标题1",
 "content": "新的内容"
}
`,
			wantCode: 200,
			wantRes:  Result{Data: float64(123)},
		},
		{
			name: "已有帖子并发表失败",
			mock: func(ctrl *gomock.Controller) (service.ArticleService, service2.InteractiveService) {
				svc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					//Id:      123,
					Title:   "新的标题1",
					Content: "新的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("系统错误"))

				return svc,intrSvc
			},
			//reqBody:  `{"id": 123,"title": "新的标题1","content":"新的内容"}`,
			reqBody: `
{
 "title": "新的标题1",
 "content": "新的内容"
}
`,
			wantCode: 200,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误"},
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造 handler
			svc ,intrSvc:= tc.mock(ctrl)
			hdl := NewArticleHandler(svc, loggerDefine.NewNopLogger(),intrSvc)

			// 准备服务器，注册路由
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UserClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRoutes(server)

			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewReader([]byte(tc.reqBody)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}
			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}