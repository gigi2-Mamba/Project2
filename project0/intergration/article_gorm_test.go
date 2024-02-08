package intergration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"project0/intergration/startup"
	"project0/internal/repository/dao"
	"project0/internal/web/ijwt"
	"testing"
)

type Result2[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}
func (s *ArticleHandlerSuite) SetupSuite() {
	s.db = startup.InitDB()
	//s.server = startup.InitWebServerJ()
	hdl := startup.InitArticleHandler(dao.NewArticleGROMDAO(s.db))
	// 只是默认启动了一个有默认为空的logger和recover的gin engine  default技巧，default基本都是写死，写空
	server := gin.Default()
	// 只是模拟token
	server.Use(func(context *gin.Context) {
		context.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoutes(server)
	s.server = server
}
func (s *ArticleHandlerSuite) TearDownTest() {
	//s.db.Exec("truncate table `articles`")
    //s.db.Exec("truncate table `published_articles`")
	//assert.NoError(s.T(), err)
}

type Article0 struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *ArticleHandlerSuite) TestArticle_Publish() {
	t := s.T()

	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		req   Article0

		// 预期响应
		wantCode   int
		wantResult2 Result2[int64]
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art dao.Article
				s.db.Where("author_id = ?", 123).First(&art)
				assert.Equal(t, "hello，你好", art.Title)
				assert.Equal(t, "随便试试", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
				assert.Equal(t, uint8(2), art.Status)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				var publishedArt dao.PublishedArticle
				s.db.Where("author_id = ?", 123).First(&publishedArt)
				assert.Equal(t, "hello，你好", publishedArt.Title)
				assert.Equal(t, "随便试试", publishedArt.Content)
				assert.Equal(t, int64(123), publishedArt.AuthorId)
				assert.Equal(t, uint8(2), publishedArt.Status)
				assert.True(t, publishedArt.Ctime > 0)
				assert.True(t, publishedArt.Utime > 0)
			},
			req: Article0{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: 200,
			wantResult2: Result2[int64]{
				Data: 1,
			},
		},

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var res Result2[int64]
			//err = json.Unmarshal(recorder.Body.Bytes(), &Result2)
			json.NewDecoder(recorder.Body).Decode(&res)
			// 当一个变量已经被assert过一次了
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult2, res)
			tc.after(t)
		})
	}
}

//借助suite包装了t，写出区别于传统的测试方法, test edit是全员通过的
//func (s *ArticleHandlerSuite) TestEdit() {
//	// 把t构建出来
//	t := s.T()
//
//	testCases := []struct {
//		name   string
//		before func(t *testing.T)
//		after  func(t *testing.T)
//
//		// 预期输入,前端传入，
//		art Article0
//		//
//		wantCode int
//		wantRes  Result2[int64]
//	}{
//
//		{
//			name:   "新建帖子",
//			before: func(t *testing.T) {},
//			after: func(t *testing.T) {
//				//验证数据库是否存储了
//				var art dao.Article
//				// 和go-redis一样查询结果链式取error   链式
//				err := s.db.Where("author_id=?", 123).Find(&art).Error
//				assert.NoError(t, err)
//				// sql执行完毕，验证其他字段
//				assert.True(t, art.Utime > 0)
//				assert.True(t, art.Ctime > 0)
//				art.Ctime = 0
//				art.Utime = 0
//
//				assert.Equal(t, dao.Article{
//					Id:      1,
//					Title:   "标题",
//					Content: "内容",
//					//Status:   1,
//					Status: domain.ArticleStatusUnPublished.ToUint8(),
//					//Fucker: 222,
//					AuthorId: 123}, art)
//			},
//			art: Article0{
//				//Id: 2,// 有id就是修改
//				Title:   "标题",
//				Content: "内容",
//			},
//			wantCode: http.StatusOK,
//			wantRes: Result2[int64]{
//				Data: 1,
//			},
//		},
//		//
//		{
//			name: "修改帖子",
//			before: func(t *testing.T) {
//				// 假装数据库已经有这个帖子
//				err := s.db.Create(&dao.Article{
//					Id:       1234,
//					Title:    "我的标题",
//					Content:  "我的内容",
//					AuthorId: 123,
//
//					Ctime:  456,
//					Status: 2,
//					Utime:  789,
//				}).Error
//				assert.NoError(t, err)
//			},
//			after: func(t *testing.T) {
//				// 你要验证，保存到了数据库里面
//				var art dao.Article
//				err := s.db.Where("id=?", 1234).
//					First(&art).Error
//				assert.NoError(t, err)
//				assert.True(t, art.Utime > 789)
//				art.Utime = 0
//				assert.Equal(t, dao.Article{
//					Id:       1234,
//					Title:    "新的标题",
//					Content:  "新的内容",
//					AuthorId: 123,
//					Status:   1,
//					Ctime:    456,
//				}, art)
//			},
//			art: Article0{
//				Id:      1234,
//				Title:   "新的标题",
//				Content: "新的内容",
//			},
//			wantCode: http.StatusOK,
//			wantRes: Result2[int64]{
//				// 我希望你的 ID 是 11
//				Data: 1234,
//			},
//		},
//		{
//			name: "修改帖子-别人的帖子",
//			before: func(t *testing.T) {
//				// 假装数据库已经有这个帖子
//				err := s.db.Create(&dao.Article{
//					Id:      3,
//					Title:   "题",
//					Content: "容",
//					// 模拟别人
//					Status:   1,
//					AuthorId: 789,
//					Ctime:    456,
//					Utime:    789,
//				}).Error
//				assert.NoError(t, err)
//				//
//
//			},
//			after: func(t *testing.T) {
//				// 你要验证，保存到了数据库里面
//				var art dao.Article
//				err := s.db.Where("id=?", 3).
//					First(&art).Error
//				assert.NoError(t, err)
//				assert.Equal(t, dao.Article{
//					Id:       3,
//					Title:    "题",
//					Content:  "容",
//					AuthorId: 789,
//					Status:   1,
//					Ctime:    456,
//					Utime:    789,
//				}, art)
//			},
//			art: Article0{
//				Id:      3,
//				Title:   "新",
//				Content: "新",
//			},
//			wantCode: http.StatusOK,
//			wantRes: Result2[int64]{
//				Msg: "系统错误",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.before(t)
//			//defer tc.after(t)
//			//defer func() {
//			//	// TRUNCATE
//			//}()
//			reqBody, err := json.Marshal(tc.art)
//			assert.NoError(t, err)
//			// 准备Req和记录的 recorder
//			req, err := http.NewRequest(http.MethodPost,
//				"/articles/edit",
//				bytes.NewReader(reqBody))
//			req.Header.Set("Content-Type", "application/json")
//			assert.NoError(t, err)
//			recorder := httptest.NewRecorder()
//			// 执行
//			s.server.ServeHTTP(recorder, req)
//			// 断言结果
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			if tc.wantCode != http.StatusOK {
//				return
//			}
//			var res Result2[int64]
//			err = json.NewDecoder(recorder.Body).Decode(&res)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.wantRes, res)
//			tc.after(t)
//		})
//	}
//}
