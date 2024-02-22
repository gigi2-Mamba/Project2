package job

import (
	"context"
	"project0/internal/service"
	"time"
)

/*
User: society-programmer on
Date: 2024/2/21  周三
Time: 19:34
*/

//
// 大写注释会有什么效果。    NewRankingJob
type  RankingJob struct {
	svc service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	return &RankingJob{svc: svc, timeout: timeout}
}

func (r *RankingJob) Name() string {

	return "ranking"
}

func (r *RankingJob) Run() error {
	//想想要run这个job 要搞些什么？   启动要控制时间
	ctx,cancel := context.WithTimeout(context.Background(),r.timeout)
	defer cancel()
	return r.svc.gTopN(ctx)
}




