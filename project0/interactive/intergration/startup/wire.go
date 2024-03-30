//go:build wireinject
package startup

import (
	"github.com/google/wire"
	"project0/interactive/grpc"
	repository2 "project0/interactive/repository"
	cache2 "project0/interactive/repository/cache"
	dao2 "project0/interactive/repository/dao"
	service2 "project0/interactive/service"
)

/*
User: society-programmer
Date: 2024/3/7  周四
Time: 18:14
*/
var thirdPartySet = wire.NewSet(
	InitDB, InitRedis, InitLogger,
	//ioc.InitRedisLimiter, ioc.NewSMSS,
	//InitSaramaClient, InitSyncProducer,
)

// grpc的wire关键就是
func InitInteractiveService() *grpc.InteractiveServiceServer {
	//wire.Build(thirdPartySet, interactiveSvcSet)
	wire.Build(thirdPartySet, interactiveSvcSet,grpc.NewInteractiveServiceServer)
	// grpc关键是返回 相应的服务的服务端实例
	return new(grpc.InteractiveServiceServer)
}


var interactiveSvcSet = wire.NewSet(
	dao2.NewInteractiveGORMDAO,
	cache2.NewInteractiveCache,
	repository2.NewCacheInteractiveRepository,
	service2.NewInteractiveService,
)


