//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	repository2 "project0/interactive/repository"
	cache2 "project0/interactive/repository/cache"
	dao2 "project0/interactive/repository/dao"
	service2 "project0/interactive/service"
	"project0/internal/events/article"
	"project0/internal/repository"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"project0/internal/service"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/ioc"
)

var thirdPartySet = wire.NewSet(
	InitDB, InitRedis, InitLogger, ioc.InitRedisLimiter, ioc.NewSMSS,
	InitSaramaClient, InitSyncProducer,
)

var userSvcProvider  = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewCacheUserRepository,
	service.NewUserService,
	)

var articleSvcProvider = wire.NewSet(
	article.NewSaramaSyncProducer,
	cache.NewArticleRedisCache,
	dao.NewArticleGROMDAO,
	//InitSyncProducer,
	repository.NewCacheArticleRepository,
	service.NewArticleService)

var interactiveSvcSet = wire.NewSet(
	dao2.NewInteractiveGORMDAO,
	cache2.NewInteractiveCache,
	repository2.NewCacheInteractiveRepository,
	service2.NewInteractiveService,
)
// 首要的main先初始化webServer
func InitWebServerJ() *gin.Engine {
	wire.Build(
		//第三方依赖，组装最基本单元
		thirdPartySet,
		userSvcProvider,
		articleSvcProvider,
		interactiveSvcSet,
		//dao.NewArticleGROMDAO,
		// cache
	 cache.NewCodeCache,
		// Repository
		// repository.NewUserRepository

		 repository.NewCodeRepository,

		// Service
		 service.NewCodeService, InitWechatService,
		ioc.InitSMSService,
		//web

		web.NewUserHandler,
		web.NewOAuth2Handler,
		ijwt.NewRedisJWTHandler,
		//InitArticleHandler,
		web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()

}

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdPartySet, interactiveSvcSet)
	return service2.NewInteractiveService(nil)
}


func InitArticleHandler(dao dao.ArticleDao) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
        userSvcProvider,
		interactiveSvcSet,

		repository.NewCacheArticleRepository,
		cache.NewArticleRedisCache,
		article.NewSaramaSyncProducer,
		service.NewArticleService,

		web.NewArticleHandler)
	return &web.ArticleHandler{}
}

