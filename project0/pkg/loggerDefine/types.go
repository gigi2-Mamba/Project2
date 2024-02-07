package loggerDefine

type Logger interface {
	Info(msg string,args...any)
	Warn(msg string,args...any)
	Debug(msg string,args...any)
	Error(msg string,args...any)
}

// 结构体约束
type LoggerV1 interface {
	Info(msg string,args...Field)
	Warn(msg string,args...Field)
	Debug(msg string,args...Field)
	Error(msg string,args...Field)
}
func example() {
	var  l Logger
	l.Info("example ","for test",3)
}
// 结构化打出



// 它要去 args 必须是偶数，并且是以 key1,value1,key2,value2 的形式传递

type LoggerV3 interface {
	Info(msg string,args...any)
	Warn(msg string,args...any)
	Debug(msg string,args...any)
	Error(msg string,args...any)
}

