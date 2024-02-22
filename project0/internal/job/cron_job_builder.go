package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"project0/pkg/loggerDefine"
	"strconv"
	"time"
)
/*
User: society-programmer on
Date: 2024/2/21  周三
Time: 19:53
*/

/**
这里做cron job仅仅只是为了监控？   是为了定时。

 */

// 构造cronjob    怎么构造
// 涉及到builder，就要有build方法
type CronJobBuilder struct {
     //既然要调用外界的东西怎么可能不会发生错误，有错误就需要log
	l loggerDefine.LoggerV1
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(l loggerDefine.LoggerV1,opt prometheus.SummaryOpts) *CronJobBuilder {
	//准备return需要的字段,这些标签名是对应到哪里的。 这可难搞啊
	prometheus.NewSummaryVec(opt,[]string{"job","success"})
	return &CronJobBuilder{l: l,
		}
}


// 利用cron去调用自己的job  这就是cron的基本灵魂
func (f *CronJobBuilder) Build(job Job) cron.Job {
     // cron原生的wrapper。 明说太死板. 在我看来是因为 原生的cron.FuncJob  只提供了一个方法。所以
	// 使用自定义的来实现cron.job。 更容易扩展。
    //return cron.FuncJob(func() {
	//})
	//为了更好的控制性，另立门户
	name := job.Name()
	return   cronJobAdapterFunc(func() {
        //为了返回一个prometheus的job
		//就是为了添加监控
		//卡壳   21.38
		start := time.Now()
		// 写一些开发日志表名启动了,关于自己研发的logger,"日志描述“ - ”详情都是需要单位对应来调用方法
		//例如string  int64  各种基本的内置变量标识符
		f.l.Debug("开始运行",loggerDefine.String("name",name))
		err := job.Run()
		if err != nil {
			//有错误肯定是 非空判断做日志好排查
			f.l.Error("执行失败",loggerDefine.Error(err),
				loggerDefine.String("name",name))
		}
		duration := time.Since(start)
		// label value 就是对应初始化vec的 label name
		f.vector.WithLabelValues(name,strconv.FormatBool(err == nil)).Observe(float64(duration.Milliseconds()))
	})

}


//实现cron的job 接口  cron.Job
type cronJobAdapterFunc func()

// 为了启动自身。   链式调用。
func (c cronJobAdapterFunc) Run()  {
	c()
}

