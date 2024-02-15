package article

import (
	"context"
	"github.com/IBM/sarama"
	"project0/internal/repository"
	"project0/pkg/loggerDefine"
	"project0/pkg/saramax"
	"time"
)

/*
Created by payden-programmer on 2024/2/14.
*/
// consumegrouphandler 的消费行为是   consumeClaim 也就是一个consume方法

// 事件相关实例，消费阅读事件的东西
type InteractiveReadEventConsumer struct {
	repo repository.InteractiveRepository
	client sarama.Client
	l loggerDefine.LoggerV1
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository, client sarama.Client, l loggerDefine.LoggerV1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
     cg,err :=sarama.NewConsumerGroupFromClient("article_read",i.client)
	//cg,err :=sarama.NewConsumerGroup(i.client)
	if err != nil {
		return err
	}
	
	go func() {
		//ctx,cancel := context.
		er := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](i.Consume, i.l))
		if er != nil {
			i.l.Error("退出消费",loggerDefine.Error(er),
				)
		}
	}()

	return err
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,evt  ReadEvent) error {
    ctx,cancel := context.WithTimeout(context.Background(),time.Second *2)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "article", evt.Aid)
}




