package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

/*
Created by payden-programmer on 2024/2/16.
*/

type Builder struct {
	Namespace string
	Subsystem string
	Name string
	InstanceId string
	Help string
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc  {
      //pattern 是指命中的路由
	labels :=[]string{"method","pattern","status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,
		Name: b.Name + "_resp_time",  //这三个字段不能有_以外的符号
		ConstLabels: map[string]string{
			"instance_id" : b.InstanceId,
		},
		// 百分比情况-误差率  50%-1%   依次类推
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	//必须注入，注册这个衡量指标。
    prometheus.MustRegister(vector)
	//有空要深究一下ctx.Next()  现在留下一个小疑问，什么情况下自己写的gin.context需要调用ctx.Next() TBC 2/21 21:52
	return func(ctx *gin.Context) {
        start := time.Now()

		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			//fullpath 是路由full path 不是整个url
			pattern := ctx.FullPath() // 可以通过上下文获取请求路径
			status := ctx.Writer.Status()
			//此时为NewVector时传递标签的值，   实例vetctor.WitLabelValues(...).Observe(毫秒数）
			vector.WithLabelValues(method,pattern,strconv.Itoa(status)).Observe(float64(duration))
		}()
        ctx.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc  {

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Name: b.Name + "_active_request",
		Help: b.Help,
		ConstLabels: map[string]string{
			"instance_id" : b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}

}
