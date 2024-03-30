package startup

import (
	"project0/internal/service/oauth2/wechat"
	"project0/pkg/loggerDefine"
)

func InitWechatService(l loggerDefine.LoggerV1) wechat.Service {

	return wechat.NewService("appID", "appSecret", l)
}
