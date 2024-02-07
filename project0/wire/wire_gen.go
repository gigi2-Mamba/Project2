// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"project0/internal/repository"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"project0/internal/service"
	"project0/internal/service/sms/failover"
	"project0/internal/web"
	"project0/internal/web/ijwt"
	"project0/ioc"
)

// Injectors from wire.go:

// 首要的main先初始化webServer
func InitWebServerJ() *gin.Engine {
	cmdable := ioc.InitRedis()
	handler := ijwt.NewRedisJWTHandler(cmdable)
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDao := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	v2 := ioc.NewSMSS()
	limiter := ioc.InitRedisLimiter(cmdable)
	smsService := ioc.InitSMSService(v2, limiter)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService(loggerV1)
	oAuth2Handler := web.NewOAuth2Handler(wechatService, userService, handler)
	articleDao := dao.NewArticleGROMDAO(db)
	articleRepository := repository.NewCacheArticleRepository(articleDao)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2Handler, articleHandler)
	return engine
}

func InitResponseTimeFailover() *failover.ResponseTimeFailover {
	v := ioc.NewSMSS()
	cmdable := ioc.InitRedis()
	limiter := ioc.InitRedisLimiter(cmdable)
	responseTimeFailover := ioc.InitFailoverService(v, limiter)
	return responseTimeFailover
}
