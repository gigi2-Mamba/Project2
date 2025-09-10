package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func initViperV1() error {
	viper.SetConfigFile("./config/dev.yaml") // 直接指定文件位置

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	log.Println("println viper ", viper.Get("test.value1"))
	log.Println("println viper grpc server addr ", viper.Get("grpc.server.addr"))
	return nil
}

// 依赖goland的configuration
func initViperDiffEnv() error {
	cfgFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	log.Println("println viper ", viper.Get("test.value1"))
	return nil
}

func initViperWatchRemote() error {
	//path   project0  / user define webook
	err := viper.AddRemoteProvider("etcd3", "http://localhost:12379", "/webook")
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置中心发生变更")
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	// 这个不合理啊   change再打印，何必隔段时间就打印呢？
	go func() {
		for {
			err = viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			//log.Println("watch ",viper.GetString("test.value1"))
			time.Sleep(10e9)
		}
	}()
	return nil
}

func initViperWatch() error {
	cfgFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	//viper.Set
	//v := viper.New()
	//viper.Debug()
	//viper.M

	viper.WatchConfig()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	log.Println("HERE CHAGNE: ", viper.GetString("test.value1"))
	//})

	log.Println("println viper ", viper.Get("test.value1"))
	return nil
}
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 用自定义的logger代替全局logger
	zap.ReplaceGlobals(logger)
}
func InitPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

// 单一服务版本
//func main() {
//
//	err := initViperDiffEnv()
//
//	if err != nil {
//		log.Println("panic in readconfig")
//		panic(err)
//	}
//
//	app := InitApp()
//	InitPrometheus()
//	for _, c := range app.consumers {
//		err := c.Start()
//		if err != nil {
//			panic(err)
//		}
//	}
//	//new Grpc server
//	server := grpc.NewServer()
//	intrv1.RegisterInteractiveServiceServer(server,app.server)
//	l,err := net.Listen("tcp",":8090")
//	if err != nil {
//		panic(err)
//	}
//	 server.Serve(l)
//
//}

// 常规的多服务版本，更灵活，更简洁
func main() {

	err := initViperV1()
	fmt.Println("123")

	if err != nil {
		log.Println("panic in read config")
		panic(err)
	}

	app := InitApp()
	fmt.Println("i'm ok app")
	InitPrometheus()
	fmt.Println("i'm ok1")
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("i'm ok2")
	err = app.server.Serve()
	if err != nil {
		panic(err)
	}
	fmt.Println("i'm ok3")

}
