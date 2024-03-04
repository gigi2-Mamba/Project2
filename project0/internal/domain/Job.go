package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

/*
User: society-programmer
Date: 2024/2/23  周五
Time: 16:00
*/


type Job struct {
	Id int64
	//同一executor 有不同任务，name来区分
	Name string
	Executor string
	//cron 表达式
	Expression string
    CancelFunc  func()

}

func (j *Job) NextTime() time.Time  {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}