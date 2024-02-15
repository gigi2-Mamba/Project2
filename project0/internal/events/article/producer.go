package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

/*
Created by payden-programmer$ on 2024/2/14$.
*/

/*
创建一个领域事件，增加阅读数
 */
//事件实体

const TopicReadEvent = "article_read"
type ReadEvent struct {
	Aid int64
	Uid int64
}
//事件方法抽象也就是事件接口
type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

// 事件实例
type SaramaSyncProducer struct {
    producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}
// produce 主要为了发送信息（SendMessage）
func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})

	return err
}

