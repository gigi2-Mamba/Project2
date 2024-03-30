package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"project0/internal/service"
	"project0/pkg/loggerDefine"
	"sync"
	"time"
)

/*
User: society-programmer on
Date: 2024/2/21  周三
Time: 19:34
*/

// 大写注释会有什么效果。    NewRankingJob
type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
	client  *rlock.Client
	key     string
	l       loggerDefine.LoggerV1
	// redis的lock
	lock *rlock.Lock
	// 报出redis的lock，
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RankingService, timeout time.Duration, l loggerDefine.LoggerV1, client *rlock.Client) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		key:       "job:ranking",
		l:         l,
		localLock: &sync.Mutex{},
		client:    client,
	}
}

func (r *RankingJob) Name() string {

	return "ranking"
}

// 保证全局只有一个实例调用任务。
func (r *RankingJob) Run() error {
	r.localLock.Lock()
	//本地锁保护分布式锁
	lock := r.lock
	if lock == nil {
		//证明没有实例获取到了分布式锁
		//去获取分布式锁的上下文
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 200,
			Max:      3,
		},
			//重试允许的最大时间
			time.Second)

		if err != nil {
			r.localLock.Unlock()
			// 这时候代表有其他实例获取到所以不管了
			r.l.Warn("ranking job 获取分布式锁失败", loggerDefine.Error(err))
			return nil
		}
		r.lock = lock
		r.localLock.Unlock()
		//自动续约
		go func() {
			// 租约间隔和续约时间
			//log.Println("timeout ",r.timeout)
			//log.Println("non-positive ",r.timeout * (4/5))
			er := lock.AutoRefresh(time.Second*315, r.timeout)

			if er != nil {
				// 续约失败了
				// 你也没办法中断当下正在调度的热榜计算（如果有）
				r.localLock.Lock()
				r.lock = nil
				r.localLock.Unlock()
			}
		}()

	}
	//到这里就是拿到锁了，开始运行job
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

//任务只在一个节点上运行

//func (r *RankingJob) Run() error {
//	//这是在job上引入分布式锁。
//	ctx,cancel := context.WithTimeout(context.Background(),time.Second * 4)
//	defer cancel()
//	lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
//		Interval: 500 * time.Millisecond,
//		Max:      3,
//		//重试的超时
//	}, time.Second*1)
//	if err != nil {
//		return err
//	}
//	// 释放锁
//	defer func() {
//		ctx,cancel := context.WithTimeout(context.Background(),time.Second * 1)
//		defer cancel()
//		er := lock.Unlock(ctx)
//		if er != nil {
//			r.l.Error("ranking_job 释放锁失败",loggerDefine.Error(er))
//			//unlock失败，不做其他处理。 因为锁本身就有过期时间。
//		}
//	}()
//
//	//想想要run这个job 要搞些什么？   启动要控制时间
//	ctx,cancel = context.WithTimeout(context.Background(),r.timeout)
//	defer cancel()
//	return r.svc.TopN(ctx)
//}
