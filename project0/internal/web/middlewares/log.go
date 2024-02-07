package middlewares

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)


//这个东西不知道什么时候写的
/*
构造一个日志中间件集成到gin当中

*/
// 日志中间件实例
type LogMiddlewareBuilder struct {
	// 这里复习一下，希望这个实例可以控制什么? 通用方法在结构体可以引入bool 来控制，直接的方法，或者匿名的方法
	// 这个方法的用处是，用context可以做链路日志
	logFn func(ctx context.Context, l AccessLog)
	// 引入bool增强控制,而怎么控制bool的值,链式调用
	allowReqBody  bool
	allowRespBody bool
}

func NewLogMiddlewareBuilder(logFn func(ctx context.Context, l AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{logFn: logFn}
}



func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder  {
	l.allowReqBody = true

	return l

}

func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder  {
	l.allowRespBody = true

	return l

}

// 辅助结构体，你希望你要记录的日志都有一些什么东西
type AccessLog struct {
	//直接抄gin的普遍实现，也是记录http请求
	Path     string    `json:"path"`
	Method   string    `json:"method"`
	ReqBody  string    `json:"reqBody"`  //
	RespBody string    `json:"respBody"` // 要设置
	Status   int       `json:"status"`
	Duration time.Duration `json:"duration"`
}

// 实际的业务逻辑在这里

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 接入系统的请求，直接说网络请求，都要考虑防御攻击
		path := c.Request.URL.Path // 请求的请求地址的路径
		if len(path) > 1024 {
			path = path[:1024]
		}

		method := c.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}
		if l.allowReqBody {
			// GetRawData is get the stream data, return the []byte
			// request body 是stream对象，只能读一次
			body, _ := c.GetRawData()

			if len(body) > 2048 {
				al.ReqBody = string(body[:2048])
			} else {
				al.ReqBody = string(body)
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		}
		start := time.Now()

		if l.allowRespBody {
			c.Writer = &responseWriter{
				al: &al,
				ResponseWriter: c.Writer,
			}
		}

		defer func() {
            al.Duration = time.Since(start)
			l.logFn(c,al)
		}()
		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)

	return w.ResponseWriter.Write(data)

}

func (w *responseWriter) WriteHeader(statusCode int) () {
	w.al.Status = statusCode

	 w.ResponseWriter.WriteHeader(statusCode)
}
