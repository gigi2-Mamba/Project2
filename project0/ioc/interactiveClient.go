package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	intrv1 "project0/api/proto/gen/api/proto/intr/v1"
	"project0/interactive/service"
	"project0/internal/client"
)

/*
User: society-programmer
Date: 2024/3/9  周六
Time: 7:08
*/

func InitIntrClient(svc service.InteractiveService)  intrv1.InteractiveServiceClient {
	// 构造依赖grpc interactive svc & local svc
	type config struct {
		Addr string `yaml:"addr"`
		Secure bool
		Threshold int32
	}

	var cfg config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	// 这种写法实则是为了对开发和线上库做分离。 真要上线得搞个else
	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts,grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	// 连接还是连接池   client connection
	log.Println(cfg.Addr)
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		// 不至于把
		panic(err)
	}
    //得到grpc的客户端
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewLocalInteractiveServiceAdapter(svc)

	interactiveClient := client.NewInteractiveClient(remote, local)
	
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = config{}
		err := viper.UnmarshalKey("grpc.client.intr",&cfg)
		log.Println("threshold : ",cfg.Threshold)
		if err != nil {
			panic(err)
		}
		interactiveClient.UpdateThreshold(cfg.Threshold)
	})

	return interactiveClient

}