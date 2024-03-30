package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"log"
	"project0/internal/job"
	"project0/internal/service"
	"project0/pkg/loggerDefine"
	"time"
)

/*
User: society-programmer on
Date: 2024/2/21  周三
Time: 22:27
*/

func InitRlockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}

func InitRankingJob(svc service.RankingService, l loggerDefine.LoggerV1, client *rlock.Client) *job.RankingJob {
	timeout := time.Minute * 4
	return job.NewRankingJob(svc, timeout, l, client)
}

// 现在就一个job所以 使用InitJobs,多个jobq其实也不存在多个job
// 提供cron实例
func InitJobs(l loggerDefine.LoggerV1, rjob *job.RankingJob) *cron.Cron {
	//
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "society_pay",
		Subsystem: "webook",
		Name:      "cron_job",
		Help:      "定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	expr := cron.New(cron.WithSeconds())

	_, err := expr.AddJob("@every 30s", builder.Build(rjob))
	if err != nil {
		log.Println("panic err here: ", err)
		panic(err)
	}
	return expr
}
