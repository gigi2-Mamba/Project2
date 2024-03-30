package article

import (
	"context"
	"github.com/IBM/sarama"
	"project0/internal/domain"
	"project0/internal/repository"
	"project0/pkg/loggerDefine"
	"project0/pkg/saramax"
	"time"
)

/*
Created by payden-programmer on 2024/2/16.
*/

type ReadHistoryConsumer struct {
	repo   repository.ReadHistoryRepository
	client sarama.Client
	l      loggerDefine.LoggerV1
}

func NewReadHistoryConsumer(repo repository.ReadHistoryRepository, client sarama.Client, l loggerDefine.LoggerV1) *ReadHistoryConsumer {
	return &ReadHistoryConsumer{repo: repo, client: client, l: l}
}

func (r *ReadHistoryConsumer) Start() error {
	// 启动消费者，需要做一个消费者组
	cg, err := sarama.NewConsumerGroupFromClient("article_record", r.client)
	if err != nil {
		return err
	}

	go func() {
		er := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if er != nil {
			r.l.Error("退出阅读历史消费", loggerDefine.Error(er))
		}
	}()

	return nil
}

func (r *ReadHistoryConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	// 先做好超时工作
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	return r.repo.AddRecord(ctx, domain.ReadHistoryRecord{
		Biz:   "article",
		BizId: evt.Aid,
		Uid:   evt.Uid,
	})

}
