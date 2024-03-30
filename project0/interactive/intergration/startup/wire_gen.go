// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/google/wire"
	"project0/interactive/grpc"
	"project0/interactive/repository"
	"project0/interactive/repository/cache"
	"project0/interactive/repository/dao"
	"project0/interactive/service"
)

// Injectors from wire.go:

// grpc的wire关键就是
func InitInteractiveService() *grpc.InteractiveServiceServer {
	cmdable := InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable)
	db := InitDB()
	interactiveDAO := dao.NewInteractiveGORMDAO(db)
	loggerV1 := InitLogger()
	interactiveRepository := repository.NewCacheInteractiveRepository(interactiveCache, interactiveDAO, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

/*
User: society-programmer
Date: 2024/3/7  周四
Time: 18:14
*/
var thirdPartySet = wire.NewSet(
	InitDB, InitRedis, InitLogger,
)

var interactiveSvcSet = wire.NewSet(dao.NewInteractiveGORMDAO, cache.NewInteractiveCache, repository.NewCacheInteractiveRepository, service.NewInteractiveService)
