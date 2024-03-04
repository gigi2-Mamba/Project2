package service

import (
	"context"
	"project0/internal/domain"
	"project0/internal/repository"
	"project0/pkg/loggerDefine"
	"time"
)

/*
User: society-programmer
Date: 2024/2/23  周五
Time: 15:44
*/

//这是为了基于mysql实现的分布式调度平台实现的jobservice
type CronJobService interface {
	Preempt(ctx context.Context) (domain.Job,error)
	ResetNextTime(ctx context.Context, job domain.Job) error
	// 设计直觉,第一种方法。 但是还有第二种更牛逼的在领域服务中处理. 所以第一种没实现。
	//Release(ctx context.Context,job domain.Job) error
	//暴露整个增删改查的post方法，TBC

}

// service惯用伎俩，把接口变小写变实例

type cronjobService struct {
	l loggerDefine.LoggerV1
	repo repository.CronJobRepository
	//抢占到job,要续约。
	refreshInterval time.Duration
}

func newCronjobService(l loggerDefine.LoggerV1, repo repository.CronJobRepository, refreshInterval time.Duration) *cronjobService {
	return &cronjobService{
		l: l,
		repo: repo,
		refreshInterval: refreshInterval}
}

func (c *cronjobService) ResetNextTime(ctx context.Context, j domain.Job) error  {
	nextTime := j.NextTime()
	return  c.repo.UpdateNextTime(ctx,j.Id,nextTime)
}

func (c *cronjobService) Preempt(ctx context.Context) (domain.Job, error) {
    j,err := c.repo.Preempt(ctx)

	if err != nil {
		return domain.Job{}, err
	}

	ticker := time.NewTicker(c.refreshInterval)
    go func() {
		//这个for range语法他妈的是哪里来的
		for range ticker.C  {
			c.refresh(j.Id)
		}
	}()
	j.CancelFunc = func() {
		//关闭续约，避免goroutine泄漏
		ticker.Stop()
		// 释放锁，很快. 不能使用入参的ctx，入参ctx天然会过期。抢占就用不上。
		ctx,cancel := context.WithTimeout(context.Background(),time.Second *1)
		defer cancel()
        err :=c.repo.Release(ctx,j.Id)
		if err != nil {
			c.l.Error("释放job失败",
				loggerDefine.Error(err),
				loggerDefine.Int64("jid",j.Id))
		}
	}


   return j, err


}


func (c *cronjobService) refresh(id int64) {
	//本质上更新一个更新时间
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	err := c.repo.UpdateTime(ctx,id)
	if err != nil {
		c.l.Error("job续约失败",loggerDefine.Error(err),
			loggerDefine.Int64("jid",id))
	}
}
