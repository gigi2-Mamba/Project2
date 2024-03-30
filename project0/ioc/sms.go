package ioc

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"project0/internal/service/sms/failover"
	"project0/pkg/limiter"

	//"project0/internal/repository"
	//"project0/internal/service"
	"project0/internal/service/sms"
	"project0/internal/service/sms/localsms"
	smstencent "project0/internal/service/sms/tecent"
)

// 是否需要说不好
//func InitService(codeRepo *repository.CodeRepository) *service.CodeService {

func InitSMSService(smss []sms.Service, limiter limiter.Limiter) sms.Service {
	//ratelimit.NewRateLimitSMSService()
	// 裝飾器模式调用第三方
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(),limiter.NewRedisSlideWindowLimiter())
	//return localsms.NewService()
	return InitFailoverService(smss, limiter)

}

// 增加服务商短信冗余
func NewSMSS() []sms.Service {
	return []sms.Service{localsms.NewService(), localsms.NewService2()}
}
func InitFailoverService(smss []sms.Service, limiter limiter.Limiter) *failover.ResponseTimeFailover {

	key := "smscode-limiter"
	return failover.NewResponseTimeFailover(smss, limiter, 3, 2, key)
}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return smstencent.NewService(c, "1400842696", "妙影科技")
}
