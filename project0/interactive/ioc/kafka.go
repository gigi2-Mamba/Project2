package startup

import (
	"github.com/IBM/sarama"
	events2 "project0/interactive/events"
	"project0/internal/events"
)

/*
Created by payden-programmer on 2024/2/14.
*/

func InitSaramaClient() sarama.Client {
	//又复习了一遍viper 走配置。  当到了ioc走配置意味是该有的整体配置已经解析出来。
	// 接下来是按需取用。 viper.Unmarshal()   or  viper.Unmarshalkey(key,cfg)
	// 地址配置，依赖配置走viper

	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"172.27.19.245:9094"}, config)
	if err != nil {
		panic(err)
	}
	return client
}

//func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
//	client, err := sarama.NewSyncProducerFromClient(c)
//	if err != nil {
//		panic(err)
//	}
//	return client
//
//}

// 因为wire很难找到同类注入，如同gin的中间件
func InitConsumers(c1 *events2.InteractiveReadEventConsumer) []events.Consumer {

	return []events.Consumer{c1}
}
