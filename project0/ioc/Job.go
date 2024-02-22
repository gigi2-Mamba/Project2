package ioc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
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


func InitRankingJob(svc service.RankingService) *job.RankingJob{
	timeout :=  time.Minute * 6
	return job.NewRankingJob(svc,timeout)
}

//现在就一个job所以 使用InitJobs
// 提供cron实例
func InitJobs(l loggerDefine.LoggerV1,rjob *job.RankingJob) *cron.Cron {
	//
	builder := job.NewCronJobBuilder(l,prometheus.SummaryOpts{
		Namespace: "society_pay",
		Subsystem: "webook",
		Name: "cron_job",
		Help: "定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	expr := cron.New(cron.WithSeconds())
	_,err := expr.AddJob("@every 1m",builder.Build(rjob))
	if err != nil {
		panic(err)
	}
	return expr
}