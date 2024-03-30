package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

/*
User: society-programmer
Date: 2024/3/8  周五
Time: 10:55
*/

// 对grpc server进行一种封装，抽取grpc服务端的通用逻辑。
type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() error  {
	l,err := net.Listen("tcp",s.Addr)
	if err != nil {
		return err
	}

	return s.Server.Serve(l)

}