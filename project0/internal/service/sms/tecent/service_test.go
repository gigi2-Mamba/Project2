package tecent

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

func TestSender(t *testing.T) {
	sid := "SMS_SECRET_ID"
	value0 := ""
	os.Setenv(sid, value0)
	skey := "SMS_SECRET_KEY"
	value1 := ""

	os.Setenv(skey, value1)
	//result, found := os.LookupEnv(key)
	//if found {
	//	fmt.Println(result)
	//} else {
	//	log.Println("system error")
	//}
	//result0, found0 := os.LookupEnv("MAKE")
	//
	//if found0 {
	//	fmt.Println(result0)
	//} else {
	//	log.Println("system error")
	//}

	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}

	s := NewService(c, "1400842696", "妙影科技")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "发送验证码",
			tplId:  "1877556",
			params: []string{"123456"},
			// 改成你的手机号码
			numbers: []string{"10086"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}

}
