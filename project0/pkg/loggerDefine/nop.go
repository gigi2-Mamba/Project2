package loggerDefine

/*
做测试用，不打印日志
*/
type NopLogger struct {
}

func (n NopLogger) Info(msg string, args ...Field) {
}

func (n NopLogger) Warn(msg string, args ...Field) {

}

func (n NopLogger) Debug(msg string, args ...Field) {

}

func (n NopLogger) Error(msg string, args ...Field) {

}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}
