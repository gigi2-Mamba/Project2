package events

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"project0/interactive/repository"
	"project0/pkg/loggerDefine"
	"project0/pkg/saramax"
	"time"
)

const TopicReadEvent = "article_read"

type ReadEvent struct {
	Aid int64
	Uid int64
}
/*
Created by payden-programmer on 2024/2/14.
*/
// consumegrouphandler 的消费行为是   consumeClaim 也就是一个consume方法

const batchSize = 10

// 事件相关实例，消费阅读事件的东西
type InteractiveReadEventConsumer struct {
	// 顾名思义   互动-阅读事件，接入互动抽象存储对象。 接入kafka解耦
	repo   repository.InteractiveRepository
	client sarama.Client
	l      loggerDefine.LoggerV1
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository, client sarama.Client, l loggerDefine.LoggerV1) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

// startBatch  v2
//func (i *InteractiveReadEventConsumer) Start() error {
//	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
//	//cg,err :=sarama.NewConsumerGroup(i.client)
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		for {
//			//ctx,cancel := context.
//			log.Println("should start  consume")
//			er := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewBatchHandler[ReadEvent](i.BatchConsume, i.l))
//			log.Println("actual consume")
//			if er != nil {
//				i.l.Error("退出消费", loggerDefine.Error(er))
//			}
//		}
//	}()
//	return err
//}

// 单个消费版本= V1
func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive_read", i.client)
	//cg,err :=sarama.NewConsumerGroup(i.client)
	if err != nil {
		return err
	}

	go func() {
		for {
			//ctx,cancel := context.
			log.Println("should start  consume")
			er := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](i.l, i.Consume))
			log.Println("actual consume")
			if er != nil {
				i.l.Error("退出消费", loggerDefine.Error(er))
			}
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	//log.Println("can't not here so far")
	return i.repo.IncrReadCnt(ctx, "article", evt.Aid)
}

//func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage, evts []ReadEvent) error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
//	defer cancel()
//	bizs := make([]string, 0, batchSize)
//	bids := make([]int64, 0, batchSize)
//	for _, evt := range evts {
//		bizs = append(bizs, "article")
//		bids = append(bids, evt.Aid)
//	}
//	//log.Println("can't not here so far")
//	return i.repo.BatchIncrReadCnt(ctx, bizs, bids)
//}
