//go:build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	 "project0/internal/repository"
	 "project0/internal/repository/cache"
	"project0/internal/repository/dao"
	 "project0/internal/service"
	"project0/internal/service/sms/failover"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/ioc"
)


var interactiveSvcSet	= wire.NewSet(
	service.NewInteractiveService,repository.NewCacheInteractiveRepository,
	cache.NewInteractiveCache,dao.NewInteractiveGORMDAO,)
// 首要的main先初始化webServer
func InitWebServerJ() *gin.Engine {
	wire.Build(
		//第三方依赖，组装最基本单元
		ioc.InitDB, ioc.InitRedis,
		//Response time trigger limiter    & 增加冗余
		ioc.InitRedisLimiter,ioc.NewSMSS,
		// DAO
		interactiveSvcSet,
		dao.NewUserDAO,dao.NewArticleGROMDAO,
		// cache
		cache.NewUserCache, cache.NewCodeCache,cache.NewArticleRedisCache,
		// Repository
		// repository.NewUserRepository
		repository.NewCacheUserRepository, repository.NewCodeRepository,repository.NewCacheArticleRepository,

		// Service
		//sms
		ioc.InitLogger,
		ioc.InitSMSService,
		// 注入这个会报错因为，没有直接入参使用到这个方法
		//ioc.InitFailoverService,// response time failover service
		service.NewUserService, service.NewCodeService,
		ioc.InitWechatService,service.NewArticleService,
		//web handler
		web.NewUserHandler,
		web.NewOAuth2Handler,
		ijwt.NewRedisJWTHandler,
		web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()

}

// 这里还可以wireinject

func InitResponseTimeFailover() *failover.ResponseTimeFailover {
	wire.Build(
		//
		ioc.InitRedis,
		//增加冗余

		ioc.NewSMSS,ioc.InitRedisLimiter,
		ioc.InitFailoverService)
	return &failover.ResponseTimeFailover{}
}
