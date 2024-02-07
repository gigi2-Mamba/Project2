package web

//import (
//	"bytes"
//	"context"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"go.uber.org/mock/gomock"
//	"golang.org/x/crypto/bcrypt"
//	"net/http"
//	"net/http/httptest"
//	"project0/internal/domain"
//	service2 "project0/internal/service"
//	svcmocks2 "project0/internal/service/mocks"
//	//"project0/internal/web/ijwt"
//	//"project0/ioc"
//	"testing"
//)
//
//// bcrypt 限制密码长度不能超过72字节
//// 思考源码真的很痛苦  这样去看源码等死吧
//func TestPassowrdEncrypt(t *testing.T) {
//	password := []byte("123456#hello")
//	encrypyted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
//	assert.NoError(t, err)
//	fmt.Println(string(encrypyted))
//	err = bcrypt.CompareHashAndPassword(encrypyted, []byte("123456#hello"))
//	assert.NoError(t, err)
//
//}
//
//func TestUserService_Signup(t *testing.T) {
//	// table driven 测试
//	testCases := []struct {
//		// 是什么测试
//		name string
//
//		//mock
//		mock func(ctrl *gomock.Controller) (service2.UserService, service2.CodeService)
//
//		// 预期输入
//		reqBuilder func(t *testing.T) *http.Request
//
//		// 预期中的输出
//		wantCode int
//		wantBody string
//	}{
//		{
//			name: "注册成功",
//			mock: func(ctrl *gomock.Controller) (service2.UserService, service2.CodeService) {
//				userSvc := svcmocks2.NewMockUserService(ctrl)
//				//codeSvc := svcmocks.NewMockCodeService(ctrl)
//				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
//					// email和phone这里我处理的有问题，还是他的版本没跟上
//					Email:    "10@qq.com",
//					Password: "123456@q",
//				}).Return(nil)
//				codeSvc := svcmocks2.NewMockCodeService(ctrl)
//				return userSvc, codeSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				jsonStr := `{"email":"10@qq.com","password":"123456@q","confirmPassword":"123456@q"}`
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(jsonStr)))
//				req.Header.Set("content-type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//			wantCode: http.StatusOK,
//			wantBody: "注册成功",
//		},
//		{
//			name: "两次密码输入不对",
//			mock: func(ctrl *gomock.Controller) (service2.UserService, service2.CodeService) {
//				userSvc := svcmocks2.NewMockUserService(ctrl)
//				//codeSvc := svcmocks.NewMockCodeService(ctrl)
//
//				codeSvc := svcmocks2.NewMockCodeService(ctrl)
//				return userSvc, codeSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				jsonStr := `{"email":"10@qq.com","password":"123456@q","confirmPassword":"123456@q3"}`
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(jsonStr)))
//				req.Header.Set("content-type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//			wantCode: http.StatusOK,
//			wantBody: "两次输入密码不一致",
//		},
//	}
//
//	// 看看你掌握了多少
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// 先生成mock实例控制器
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish() // 关闭mock 好习惯
//			// mock模拟依赖项
//			userSvc, codeSvc := tc.mock(ctrl)
//			// 依赖项注入
//			//jwthdl := ijwt.NewRedisJWTHandler(ioc.InitRedis())
//			hdl := NewUserHandler(userSvc, codeSvc)
//			// 构建gin 实例
//			server := gin.Default()
//			// 注册路由
//			hdl.RegisterRoute(server)
//			// 构造请求
//			req := tc.reqBuilder(t)
//			recorder := httptest.NewRecorder()
//			server.ServeHTTP(recorder, req)
//
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			assert.Equal(t, tc.wantBody, recorder.Body.String())
//		})
//
//	}
//}
//
//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	userSvc := svcmocks2.NewMockUserService(ctrl)
//	var u = domain.User{
//		Id:    4,
//		Email: "2@qq.com"}
//	userSvc.EXPECT().Signup(gomock.Any(), u).Return(errors.New("db 出错"))
//
//	err := userSvc.Signup(context.Background(), u)
//	t.Log(err)
//
//}
