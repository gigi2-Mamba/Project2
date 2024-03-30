package web

// 单元测试做handler
//func TestArticleHandler_Publish(t *testing.T) {
//	testCases := []struct {
//		name string
//		mock func(ctrl *gomock.Controller) service.ArticleService
//
//		reqBody  string
//		wantCode int
//		wantRes  Result
//	}{
//		{
//			name: "新建并发表成功xxxx",
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
//			name: "新建并发表成功",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Title:   "我的标题",
//					Content: "我的内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(1), nil)
//				return svc
//			},
//			reqBody:  `{ "title":"我的标题","content":"我的内容"}`,
//			wantCode: http.StatusOK,
//			//在json格式中any类型默认使用int64
//			wantRes: Result{Data: float64(1)},
//		},
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
//		//{},
//		{
//			name: "已有帖子并发表失败2",
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
//			name: "已有帖子并发表失败3",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Id:      123,
//					Title:   "我的标题",
//					Content: "我的内容2",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(123), nil)
//				return svc
//			},
//			reqBody:  `{"id": 123,"title": "我的标题","content":"我的内容2"}`,
//			wantCode: http.StatusOK,
//			//在json格式中any类型默认使用int64
//			wantRes: Result{Data: float64(123)},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish() // 确保期望被调用的都调用了  做守门人
//			svc := tc.mock(ctrl)
//			// 构造调用母体
//
//			handler := NewArticleHandler(svc, loggerDefine.NewNopLogger())
//			// 寫到这里就很关键了，如何mock构造一个http请求
//			// 准备服务器和注册路由
//			server := gin.Default()
//			server.Use(func(ctx *gin.Context) {
//				ctx.Set("user", ijwt.UserClaims{
//					Uid: 123,
//				})
//			})
//			handler.RegisterRoutes(server)
//			// 构造请求
//			// bytes.NewReader([]byte)  和  bytes.NewBufferString() 就是少了个[]byte，好像更便捷。
//			// 其实用reader是为了只读，buffer实际上可以修改。 但是在这里构造一个http请求走bufferstring省事
//			req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBufferString(tc.reqBody))
//			//通俗解释,写入请求头是为了告诉服务器，当前请求发送的数据类型是什么？  是json   内容类型
//			assert.NoError(t, err)
//			req.Header.Set("Content-Type", "application/json")
//			// 构造响应
//			recorder := httptest.NewRecorder()
//			server.ServeHTTP(recorder, req)
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			if recorder.Code != http.StatusOK {
//				log.Println("here test gogogo")
//				return
//			}
//			var res Result
//
//			err = json.NewDecoder(recorder.Body).Decode(&res)
//			assert.NoError(t, err)
//			//assert.Equal(t, tc.wantCode, recorder.Code)
//			assert.Equal(t, tc.wantRes, res)
//		})
//	}
//}
