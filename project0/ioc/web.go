package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/internal/web/middlewares"
	"project0/pkg/ginx"
	"project0/pkg/ginx/middleware/prometheus"
	"project0/pkg/loggerDefine"
	"strings"
	"time"
)

// 这样的注入web
// 先注册中间件
func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2Handler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default() // default version gin.Engine
	server.Use(mdls...)     // 用的radix tree回头回顾一下 ，先插入中间件undoubtedly
	// 后续增加了handler,就要继续补充。  hdl.RegisterRoute(server) 对gin server 注册路由
	userHdl.RegisterRoute(server)
	wechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

// 要使用基于redis的ratelimit，要加入redis.cmdable
func InitGinMiddlewares(redisClient redis.Cmdable, Hdl ijwt.Handler, l loggerDefine.LoggerV1) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "geektime_daming",
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "society_pay",
		Subsystem: "webook",
		Name:      "biz_code",
		Help:      "统计业务错误码	",
	})
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true, // 允许携带cookie
			AllowHeaders:     []string{"Content-Type", "authorization"},
			AllowOrigins:     []string{"http://localhost:3000"},
			//AllowAllOrigins: true,
			//AllowOriginFunc: func(origin string) bool {
			//		//	return  true
			//		//},
			MaxAge: 12 * time.Hour,
			//允许前端访问后端响应携带的token
			AllowOriginFunc: func(origin string) bool {
				log.Println(strings.HasPrefix(origin, "http://localhost"))
				if strings.HasPrefix(origin, "http://localhost") {
					//log.Println("may be here")
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			// 为了让前端可以拿到，后端做的 ctx.Header(key,value), 所以exposeHeaders for the header kv u set
			//jwt token是短token,refresh token 是长token
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		}), func(context *gin.Context) {
			//log.Println("跨域通过middleware")
		},
		otelgin.Middleware("webook"),
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		//ratelimit.NewBuilder(limiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 1000)).Build(),
		//middlewares.NewLogMiddlewareBuilder(func(ctx context.Context, al middlewares.AccessLog) {
		//	l.Debug("",loggerDefine.Field{Key: "req",Val: al})
		//}).AllowReqBody().AllowRespBody().Build(),
		middlewares.NewLoginJWTMiddlewareBuilder(Hdl).CheckLogin(),
	}
}
