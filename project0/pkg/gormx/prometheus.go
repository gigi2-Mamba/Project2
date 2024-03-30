package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

/*
Created by payden-programmer on 2024/2/17.
*/

// 为了进一步封装，直接实现gorm的plugin接口方便，db.use
type Callbacks struct {
	//为什么用矢量呢？
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", c.Before())
	if err != nil {
		return err
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", c.After("CREATE"))
	if err != nil {
		return err
	}

	err = db.Callback().Query().Before("*").
		Register("prometheus_query_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("*").
		Register("prometheus_query_after", c.After("QUERY"))
	if err != nil {
		return err
	}

	err = db.Callback().Query().Before("*").
		Register("prometheus_raw_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", c.After("RAW"))
	if err != nil {
		return err
	}

	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", c.After("UPDATE"))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").
		Register("prometheus_delete_after", c.After("DELETE"))
	if err != nil {
		return err
	}

	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").
		Register("prometheus_row_after", c.After("ROW"))
	return err
}

func NewCallbacks(opt prometheus.SummaryOpts) *Callbacks {
	vector := prometheus.NewSummaryVec(opt, []string{"type", "table"})
	//重点，注册，注册
	prometheus.MustRegister(vector)
	return &Callbacks{
		//想要知道的label
		vector: vector,
	}
}

func (c *Callbacks) After(typ string) func(db *gorm.DB) {

	return func(db *gorm.DB) {
		//能省一点是一点，靠推断
		val, _ := db.Get("start_time")
		start, ok := val.(time.Time)
		if ok {
			duration := time.Since(start).Milliseconds()

			c.vector.WithLabelValues(typ, db.Statement.Table).Observe(float64(duration))
		}

	}

}
func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}
