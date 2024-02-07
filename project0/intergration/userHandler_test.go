package intergration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"project0/intergration/startup"
	"project0/internal/web"
	"testing"
	"time"
)

func TestUserHandler_sendSmsCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServerJ()
	testCases := []struct {
		name string

		//before
		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wantBody web.Result
	}{
		{name: "发送验证码",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				//after need   验证验证码是否存在，是否存在存在时间目的就结束了
				key := "phone_code:bizLogin:13168896092"
				contextTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				// redis命令返回的结果比较通用
				code, err := rdb.Get(contextTimeout, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				duration, err := rdb.TTL(contextTimeout, key).Result()

				assert.NoError(t, err)
				assert.True(t, duration > time.Minute*9+50)
				//因为是模拟的还要删除
				err = rdb.Del(contextTimeout, key).Err()
				assert.NoError(t, err)

			},
			phone:    "13168896092",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// 构造req
			request, err := http.NewRequest(http.MethodPost, "/users/login_sms/send/code0",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			request.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, request)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			//log.Println("res: ",res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)

		})
	}
}
