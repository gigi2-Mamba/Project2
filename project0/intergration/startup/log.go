package startup

import "project0/pkg/loggerDefine"

func InitLogger() loggerDefine.LoggerV1 {
	return  loggerDefine.NewNopLogger()
}
