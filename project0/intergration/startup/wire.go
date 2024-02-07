//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	repository2 "project0/internal/repository"
	cache2 "project0/internal/repository/cache"
	"project0/internal/repository/dao"
	service2 "project0/internal/service"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/ioc"
)

var thirdPartySet = wire.NewSet(
	InitDB, InitRedis,InitLogger,ioc.InitRedisLimiter,ioc.NewSMSS,
	)
// 首要的main先初始化webServer
func InitWebServerJ() *gin.Engine {
	wire.Build(
		//第三方依赖，组装最基本单元
		thirdPartySet,
		// DAO
		dao.NewUserDAO,dao.NewArticleGROMDAO,
		// cache
		cache2.NewUserCache, cache2.NewCodeCache,
		// Repository
		// repository.NewUserRepository
		repository2.NewCacheUserRepository, repository2.NewCodeRepository,repository2.NewCacheArticleRepository,
		// Service
		service2.NewUserService, service2.NewCodeService,service2.NewArticleService,InitWechatService,
		ioc.InitSMSService,
		//web

		web.NewUserHandler,
		web.NewOAuth2Handler,
        ijwt.NewRedisJWTHandler,
        web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()

}

func InitArticleHandler(dao dao.ArticleDao) *web.ArticleHandler  {
	wire.Build(
		thirdPartySet,
		service2.NewArticleService,
		repository2.NewCacheArticleRepository,
		web.NewArticleHandler,)
	return  &web.ArticleHandler{}
}
