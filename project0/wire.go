//go:build wireinject

package main

import (
	"github.com/google/wire"
	"project0/interactive/events"
	repository2 "project0/interactive/repository"
	cache2 "project0/interactive/repository/cache"
	dao2 "project0/interactive/repository/dao"
	service2 "project0/interactive/service"
	"project0/internal/events/article"
	"project0/internal/repository"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"project0/internal/service"
	"project0/internal/service/sms/failover"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/ioc"
)

var interactiveSvcSet = wire.NewSet(
	service2.NewInteractiveService, repository2.NewCacheInteractiveRepository,
	cache2.NewInteractiveCache, dao2.NewInteractiveGORMDAO,)

var rankSvcSet = wire.NewSet(
	cache.NewRankingRedisCache, repository.NewCacheRankingRepository,
	service.NewBatchRankingService)

// 首要的main先初始化webServer
func InitWebServerJ() *App {
	wire.Build(
		//第三方依赖，组装最基本单元
		ioc.InitDB, ioc.InitRedis,
		//Response time trigger limiter    & 增加冗余
		ioc.InitRedisLimiter, ioc.NewSMSS,
		ioc.InitSaramaClient, ioc.InitSyncProducer,
		ioc.InitConsumers, ioc.InitRankingJob, ioc.InitJobs,
		ioc.InitRlockClient,
		// DAO
		interactiveSvcSet, rankSvcSet,
		ioc.InitIntrClient,
		article.NewSaramaSyncProducer, events.NewInteractiveReadEventConsumer, article.NewReadHistoryConsumer,
		dao.NewUserDAO, dao.NewArticleGROMDAO, dao.NewHistoryGORMDAO,
		// cache
		cache.NewUserCache, cache.NewCodeCache, cache.NewArticleRedisCache,
		// Repository
		// repository.NewUserRepository
		repository.NewCacheUserRepository, repository.NewCodeRepository, repository.NewCacheArticleRepository,
		repository.NewCacheArticleHistoryRepository,
		// Service
		//sms
		ioc.InitLogger,
		ioc.InitSMSService,
		// 注入这个会报错因为，没有直接入参使用到这个方法
		//ioc.InitFailoverService,// response time failover service
		service.NewUserService, service.NewCodeService,
		ioc.InitWechatService, service.NewArticleService,
		//web handler
		web.NewUserHandler,
		web.NewOAuth2Handler,
		ijwt.NewRedisJWTHandler,
		web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
		// wire大结局，一个app代表整个应用,这个也值得回头一看
		wire.Struct(new(App), "*"),
	)
	return new(App)

}

// 这里还可以wireinject

func InitResponseTimeFailover() *failover.ResponseTimeFailover {
	wire.Build(
		//
		ioc.InitRedis,
		//增加冗余

		ioc.NewSMSS, ioc.InitRedisLimiter,
		ioc.InitFailoverService)
	return &failover.ResponseTimeFailover{}
}
