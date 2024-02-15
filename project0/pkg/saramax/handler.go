package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"

	"project0/pkg/loggerDefine"
)

/*
Created by payden-programmer on 2024/2/14.
*/

// 为什么要写这个handler   这个handler是单消费版本
type Handler[T any] struct {
	// 这个方法是怎么想出来的，很自然而然需要消费者信息和事件
	fn func(msg *sarama.ConsumerMessage,evt T) error
	l loggerDefine.LoggerV1
}

// 泛型方法。 返回值是指定的T就可以
func NewHandler[T any]( l loggerDefine.LoggerV1,fn func(msg *sarama.ConsumerMessage, evt T) error) *Handler[T] {
	return &Handler[T]{fn: fn, l: l}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// 怎么消费信息？
func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 复习一下，ConsumeClaim会完成什么功能？
	//先获取消息管道
	msgs := claim.Messages()

	for msg := range msgs {
		//业务逻辑
		//直接依赖泛型约束
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化消息失败",
				loggerDefine.String("topic",msg.Topic),
				loggerDefine.Int32("partition",msg.Partition),
				loggerDefine.Int64("offset",msg.Offset),
				loggerDefine.Error(err))
			return err

		}
		err = h.fn(msg,t)
		if err != nil {
			h.l.Error("处理消息失败",
				loggerDefine.String("topic",msg.Topic),
				loggerDefine.Int32("partition",msg.Partition),
				loggerDefine.Int64("offset",msg.Offset),
				loggerDefine.Error(err))
			return err

		}
       // 标记消息被消费了 ，为什么一定要标记啊？
		session.MarkMessage(msg,"处理阅读数")
	}
	return nil

}