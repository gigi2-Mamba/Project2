package ioc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

// viper序列化后是通过map[string]interface{} 去反序列化得配置然后再调用正经得初始化方法
func InitRedis() redis.Cmdable {
	type cfg struct {
		//用yaml 很舒服 暂时不需要其他部分引用，所以使用匿名结构体
		Addr string `yaml:"addr"`
	}

	var rcfg cfg
	viper.UnmarshalKey("redis", &rcfg)
	log.Println("rcfg is : ", rcfg.Addr)
	//log.Println("whether true: ",rcfg.Addr == "")
	client := redis.NewClient(&redis.Options{
		Addr:     rcfg.Addr,
		Password: "lxjredis123",
	})
	ping := client.Ping(context.Background()).Err()
	if ping != nil {
		log.Println("ping err: ", ping)
		return nil
	}
	fmt.Println("redis connected successfully")
	return client
}
