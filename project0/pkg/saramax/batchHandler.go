package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"project0/pkg/loggerDefine"
	"time"
)

/*
Created by payden-programmer on 2024/2/15.
*/

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, events []T) error
	l  loggerDefine.LoggerV1
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// 以後可以尝试重复写多次这个方法来掌握批量消费
func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//先把消息搞出来
	msgs := claim.Messages()
	const batchSize = 10 //定义一个批量先,批量要注意超时消费
	for {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		// 然后批量处理消息
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		var done = false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					// 消息通道关闭
					cancel()
					return nil
				}
				//batch = append(batch,msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化消息失败",
						loggerDefine.String("topic", msg.Topic),
						loggerDefine.Int32("partition", msg.Partition),
						loggerDefine.Int64("offset", msg.Offset),
						loggerDefine.Error(err))
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		// 凑够了一批，然后你就处理
		err := b.fn(batch, ts)
		if err != nil {
			b.l.Error("处理消息失败",
				// 把真个 msgs 都记录下来
				loggerDefine.Error(err))
			
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}

	}

}

func NewBatchHandler[T any](fn func(msgs []*sarama.ConsumerMessage, events []T) error, l loggerDefine.LoggerV1) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}
