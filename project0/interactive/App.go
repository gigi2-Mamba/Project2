package main

import (
	"project0/internal/events"
	"project0/pkg/grpcx"
)

//单服务版本
//type App struct {
//
//	consumers []events.Consumer
//	server *grpc.InteractiveServiceServer
//
//}


type App struct {

	consumers []events.Consumer
	// 装饰器模式
	server *grpcx.Server

}
