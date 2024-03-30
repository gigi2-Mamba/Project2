//go:build wireinject
package main

import (
	"github.com/google/wire"
	"project0/interactive/events"
	"project0/interactive/grpc"
	ioc "project0/interactive/ioc"
	repository2 "project0/interactive/repository"
	cache2 "project0/interactive/repository/cache"
	dao2 "project0/interactive/repository/dao"
	service2 "project0/interactive/service"

)

/*
User: society-programmer
Date: 2024/3/8  周五
Time: 9:29
*/

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		interactiveSvcSet,
		InitInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App),"*"))
	return new(App)
}

// 依赖注入 Interactive的所有第三方依赖
var  thirdPartySet = wire.NewSet(
	ioc.InitDB,ioc.InitRedis,ioc.InitLogger,ioc.InitSaramaClient,)

var interactiveSvcSet = wire.NewSet(
	dao2.NewInteractiveGORMDAO,
	cache2.NewInteractiveCache,
	repository2.NewCacheInteractiveRepository,
	service2.NewInteractiveService,
)

func InitInteractiveServiceServer() *grpc.InteractiveServiceServer {
	//wire.Build(thirdPartySet, interactiveSvcSet)
	wire.Build(thirdPartySet, interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		)
	// grpc关键是返回 相应的服务的服务端实例
	return new(grpc.InteractiveServiceServer)
}



