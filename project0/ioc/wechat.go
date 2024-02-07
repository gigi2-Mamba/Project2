package ioc

import (
	"os"
	"project0/internal/service/oauth2/wechat"
	"project0/pkg/loggerDefine"
)

func InitWechatService(l loggerDefine.LoggerV1) wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")

	if !ok {
		appID = "KKKK"
		//panic("找不到环境变量")
	}

	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	appSecret ="DFDFDF"
	return wechat.NewService(appID, appSecret,l)
}
