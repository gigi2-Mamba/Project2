//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
	InitDB, InitRedis,InitLogger,ioc.InitRedisLimiter,ioc.NewSMSS,
	InitSaramaClient,InitSyncProducer,
	)
// 首要的main先初始化webServer
func InitWebServerJ() *gin.Engine {
	wire.Build(
		//第三方依赖，组装最基本单元
		thirdPartySet,
		// DAO
		dao.NewUserDAO,dao.NewArticleGROMDAO,
		// cache
		cache.NewUserCache, cache.NewCodeCache,
		// Repository
		// repository.NewUserRepository

		repository.NewCacheUserRepository, repository.NewCodeRepository,

		// Service
		service.NewUserService, service.NewCodeService,InitWechatService,
		ioc.InitSMSService,
		//web

		web.NewUserHandler,
		web.NewOAuth2Handler,
        ijwt.NewRedisJWTHandler,
		InitArticleHandler,
        //web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()

}

func InitArticleHandler(daoArt dao.ArticleDao) *web.ArticleHandler  {
	wire.Build(
		thirdPartySet,
		dao.NewUserDAO,
		cache.NewInteractiveCache,
		article.NewSaramaSyncProducer,
		dao.NewInteractiveGORMDAO,
		repository.NewCacheInteractiveRepository,
		service.NewArticleService,
		service.NewInteractiveService,
		repository.NewCacheArticleRepository,
		cache.NewArticleRedisCache,
		web.NewArticleHandler,)
	return  &web.ArticleHandler{}
}
