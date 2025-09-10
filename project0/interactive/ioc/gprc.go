package startup

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	grpc2 "project0/interactive/grpc"
	"project0/pkg/grpcx"
)

/*
User: society-programmer
Date: 2024/3/8  周五
Time: 11:12
*/

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	// 一般会想到,最原始的做法。 但是有更巧妙的方法。
	//intrv1.RegisterInteractiveServiceServer()

	addr := viper.GetString("grpc.server.addr")
	log.Println("addr : ", addr)
	if addr == "" { // 没有读取成功是怎么回事？
		panic(addr)
	}
	return &grpcx.Server{
		Server: s,
		// 由于单独的服务配置比较简单。  用这个GetString。 这个string是整个链路配置的全路径。直接拿到值。
		Addr: viper.GetString("grpc.server.addr"),
	}
}
