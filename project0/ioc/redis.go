package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

func InitRedis() redis.Cmdable {
     type cfg struct {
		 //用yaml 很舒服 暂时不需要其他部分引用，所以使用匿名结构体
		 addr string `yaml:"addr"`
	 }

	 var rcfg cfg
	 viper.UnmarshalKey("redis",&rcfg)
	 log.Println("rcfg is : ",rcfg.addr)
	return redis.NewClient(&redis.Options{Addr: rcfg.addr})
}
