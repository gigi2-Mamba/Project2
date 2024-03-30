package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"log"
	events2 "project0/interactive/events"
	"project0/internal/events"
	"project0/internal/events/article"
)

/*
Created by payden-programmer on 2024/2/14.
*/

func InitSaramaClient() sarama.Client {
	//又复习了一遍viper 走配置。  当到了ioc走配置意味是该有的整体配置已经解析出来。
	// 接下来是按需取用。 viper.Unmarshal()   or  viper.Unmarshalkey(key,cfg)
	// 地址配置，依赖配置走viper
	type Cfg struct {
		Addr []string
	}
	var cfg Cfg
	// 不管错误，不应该有错误？
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, config)
	if client == nil {
		log.Println("hereskdslfjdsl:", cfg.Addr)
		panic(client)
	}
	return client

}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {

	client, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return client

}

// 因为wire很难找到同类注入，如同gin的中间件
// 这里有两个consumer
func InitConsumers(c1 *events2.InteractiveReadEventConsumer, c2 *article.ReadHistoryConsumer) []events.Consumer {

	return []events.Consumer{c1, c2}
}
