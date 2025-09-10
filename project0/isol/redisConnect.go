package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		DB:       0,
		Password: "lxjredis123"})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	client.Set(context.Background(), "key1", "value1", 80e9)
}
