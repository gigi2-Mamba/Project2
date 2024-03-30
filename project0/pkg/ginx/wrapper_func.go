package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"net/http"
	"project0/pkg/loggerDefine"
	"strconv"
)

func NewLogger() *zap.Logger {
	development, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return development
}

var logger = NewLogger()
var L loggerDefine.LoggerV1 = loggerDefine.NewZapLogger(logger)

// 受制于泛型，只能用包变量
var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func WrapBody[Req any](bizFn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", loggerDefine.Error(err))
		}
		//L.Debug("输入参数",loggerDefine.Field{Key: "req",Val: req})

		res, err := bizFn(ctx, req)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", loggerDefine.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBodyAndClaims[Req any, Claims jwt.Claims](
	bizFn func(ctx *gin.Context, req Req, uc Claims) (Result, error)) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入有误", loggerDefine.Error(err))
		}
		L.Debug("输入参数", loggerDefine.Field{"req", req})
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := bizFn(ctx, req, uc)
		//简单的计数是什么情况
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("业务逻辑执行失败", loggerDefine.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
