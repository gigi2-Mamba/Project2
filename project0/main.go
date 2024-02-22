package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project0/internal/service/sms/failover"
	"project0/ioc"
	"time"
)

func init() {
	//gin.SetMode(gin.ReleaseMode)
}

func initViper() error {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	//当前工作目录下的子包config,goland的工作目录，如果不是goland，那得小心当前开发平台的工作目录
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	log.Println("println viper ",viper.Get("test.value1"))
    return nil
}

func initViperV1() error {
	viper.SetConfigFile("./config/dev.yaml")  // 直接指定文件位置

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	log.Println("println viper ",viper.Get("test.value1"))
	return nil
}

// 依赖goland的configuration
func initViperDiffEnv() error {
	cfgFile := pflag.String("config","config/dev.yaml","配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		//log.Println("读取配置有误",err)
		panic(err)
	}
	log.Println("println viper ",viper.Get("test.value1"))
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
	cfgFile := pflag.String("config","config/dev.yaml","配置文件路径")
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
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("HERE CHAGNE: ",viper.GetString("test.value1"))
	})

	log.Println("println viper ",viper.Get("test.value1"))
	return nil
}
func InitLogger()  {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 用自定义的logger代替全局logger
	zap.ReplaceGlobals(logger)
}
func InitPrometheus()  {
	go func() {
		http.Handle("/metrics",promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()

}
func main() {
	//server := wire.InitWebServerJ()	//err := mysqlInit.Init()
	//	//if err != nil {
	//	//	return
	//	//}
	//	//log.Println("mysqlInit init success")
	//	server := initWebServer()
	//err := initViperWatchRemote()
	err := initViperV1()
	InitLogger()
	tpCancel := ioc.InitOTEL()
	defer func() {
		ctx,cancel := context.WithTimeout(context.Background(),time.Second)
		defer cancel()
		tpCancel(ctx)

	}()
	if err != nil {
		log.Println("panic in readconfig")
		panic(err)
	}

	app :=InitWebServerJ()
	InitPrometheus()
	for _,c := range app.consumers {
		//log.Println("consumer is : ",c)
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	app.cron.Start()
	defer func() {
		<- app.cron.Stop().Done()
	}()

	server := app.server
	//server := wire.InitWebServerJ()
	//server.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello  k8s   部署成功了")
	//})
	//
	//server.GET("/setcookie", func(context *gin.Context) {
	//	context.SetCookie("cookietst", "123", 300, "/oauth2/wechat/getcookie", "", false, true)
	//	context.String(http.StatusOK, "设置cookie成功？")
	//})
	//initUserHdl(mysqlInit.Db, server)
	go failover.AsyncSendCode(InitResponseTimeFailover())
	server.Run(":8083")
}

// session 和 jwt可以交替使用
//func useSession(server *gin.Engine) {
//
//	login := &middlewares.LoginMiddlewareBuilder{}
//	// 先初始化 存储数据,直接存cookie作教学
//	//store := cookie.NewStore([]byte("secret"))
//	//  注意localhost  是否讲得通 ubuntu ,  两个密钥 Authentication, encryption  身份验证，数据加密，授予权限
//	// store 面向接口编程
//	// gin的redis模块不走client,直接走NewStore 来存储  tbc复杂的设计
//	store, err := redis.NewStore(16, "tcp",
//		"localhost:6379", "",
//		[]byte("oDhIbNhVlYcOtAqNvVaMlFbQrDdObWqT"),
//		[]byte("oDhIbNhVlYcOtAqNvVaMlFbQrDdObWxT"))
//	fmt.Println("store err,session err: ", err, store)
//	//handlerFunc := sessions.Sessions("ssid", store)
//	server.Use(sessions.Sessions("ssig", store), login.CheckLogin())
//	log.Println("can here?  ")
//	if err != nil {
//		panic(err)
//	}
//}
