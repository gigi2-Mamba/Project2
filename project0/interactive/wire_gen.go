// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"project0/interactive/events"
	"project0/interactive/grpc"
	"project0/interactive/ioc"
	"project0/interactive/repository"
	"project0/interactive/repository/cache"
	"project0/interactive/repository/dao"
	"project0/interactive/service"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitApp() *App {
	cmdable := startup.InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable)
	db := startup.InitDB()
	interactiveDAO := dao.NewInteractiveGORMDAO(db)
	loggerV1 := startup.InitLogger()
	interactiveRepository := repository.NewCacheInteractiveRepository(interactiveCache, interactiveDAO, loggerV1)
	client := startup.InitSaramaClient()
	interactiveReadEventConsumer := events.NewInteractiveReadEventConsumer(interactiveRepository, client, loggerV1)
	v := startup.InitConsumers(interactiveReadEventConsumer)
	interactiveServiceServer := InitInteractiveServiceServer()
	server := startup.NewGrpcxServer(interactiveServiceServer)
	app := &App{
		consumers: v,
		server:    server,
	}
	return app
}

func InitInteractiveServiceServer() *grpc.InteractiveServiceServer {
	cmdable := startup.InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable)
	db := startup.InitDB()
	interactiveDAO := dao.NewInteractiveGORMDAO(db)
	loggerV1 := startup.InitLogger()
	interactiveRepository := repository.NewCacheInteractiveRepository(interactiveCache, interactiveDAO, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

// 依赖注入 Interactive的所有第三方依赖
var thirdPartySet = wire.NewSet(startup.InitDB, startup.InitRedis, startup.InitLogger, startup.InitSaramaClient)

var interactiveSvcSet = wire.NewSet(dao.NewInteractiveGORMDAO, cache.NewInteractiveCache, repository.NewCacheInteractiveRepository, service.NewInteractiveService)
