package ioc

import (
	"go.uber.org/zap"
	"project0/pkg/loggerDefine"
)

func InitLogger() loggerDefine.LoggerV1  {
   l,err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}


   return loggerDefine.NewZapLogger(l)
}
