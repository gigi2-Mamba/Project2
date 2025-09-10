package startup

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitRedis() redis.Cmdable {
	fmt.Println("try to initial redis")
	client := redis.NewClient(&redis.Options{Addr: "localhost:16379"})
	// 最近添加
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Println("interactive redis connected err: ", err.Error())
	}
	return client
}
