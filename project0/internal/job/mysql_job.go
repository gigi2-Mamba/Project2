package job

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"project0/internal/domain"
	"project0/internal/service"
	"project0/pkg/loggerDefine"
	"time"
)

/*
User: society-programmer
Date: 2024/2/24  周六
Time: 10:28
*/
//提供执行器的名字，和执行方法
type Executor interface {
	//执行器的名字，走按名索引
	Name() string
    //ctx  全局控制 executor 要注意ctx的超时和取消
	Exec(ctx context.Context,j domain.Job) error
}

//为什么有这个本地方法的东西干春啊
// 本地执行器。
type  LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context,j domain.Job) error
}

func NewLocalFuncExecutor(funcs map[string]func(ctx context.Context, j domain.Job) error) *LocalFuncExecutor {
	return &LocalFuncExecutor{funcs: funcs}
}

//分布式调度器
type Scheduler struct {
	//定时任务服务
	svc service.CronJobService
	// 调度器也有一个时间,这个怎么好像无限复用啊  TBC
	dbTimeout time.Duration
	// 执行器列表。
	executors  map[string]Executor
	l  loggerDefine.LoggerV1
	// 限制一直抢占。保持系统稳定性。   新东西  semaphore.Weighted
	//这个东西很简便
	limiter  *semaphore.Weighted

}

func NewScheduler(svc service.CronJobService, l loggerDefine.LoggerV1) *Scheduler {
	return &Scheduler{
		svc: svc,
		l: l,
	    dbTimeout: time.Second,
	    limiter: semaphore.NewWeighted(100),
	    executors: map[string]Executor{},
	}
}
func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn,ok := l.funcs[j.Name]
	if !ok {
		return  fmt.Errorf("未注冊本地方法%s",j.Name)
	}
	return fn(ctx,j)
}


func (l *LocalFuncExecutor) RegisterExecutor(name string, fn func(ctx context.Context,j domain.Job) error )  {
	 l.funcs[name] = fn
}


func (s *Scheduler) RegisterExecutor(exec Executor)  {
	s.executors[exec.Name()] = exec
}

//调度器的调度方法
func (s *Scheduler) Schedule(ctx context.Context)   {
	for  { //不断的抢占，所以要限制
		// 先判断上下文是否有问题
		// TBC
		if  ctx.Err() != nil {
			return
		}
		err := s.limiter.Acquire(ctx,1)
		if err != nil {
			//当出现err时 就不允许调度了。所以做了一个限制。 就是一个简易令牌。
			return
		}
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		job, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			//发生错误,继续调度
			continue
		}
		//调度执行j
        exec, ok := s.executors[job.Executor]
		if !ok {
			//找不到执行器，记录一下
            s.l.Error("找不到执行器",loggerDefine.Int64("jid",job.Id),
				loggerDefine.String("executor",job.Executor))
			continue
		}
		
		//  异步执行不然会阻塞
		go func() {
			defer func() {
				s.limiter.Release(1)
				job.CancelFunc()
			}()
			er := exec.Exec(ctx, job)
			if er != nil{
				s.l.Error("执行任务失败",loggerDefine.Error(er),
					loggerDefine.Int64("jid",job.Id), )

			}
			er = s.svc.ResetNextTime(ctx,job)
			if er != nil {
				s.l.Error("重置下次执行时间失败",loggerDefine.Error(er),
					loggerDefine.Int64("jid",job.Id),)
			}
		}()

	}
}

