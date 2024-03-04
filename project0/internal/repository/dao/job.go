package dao

import (
	"context"
	"gorm.io/gorm"
	"project0/internal/domain"
	"time"
)

/*
User: society-programmer
Date: 2024/2/23  周五
Time: 16:19
*/


//先定义空壳子再说，卧槽
type Job struct {
	Id  int64
	Status int
	// 引入标准的乐观锁,利用version可以确保没有人改
	Version int

	Expression string
	Executor string
	Name string `gorm:"unique"`
    //用户额外的配置，实例id，环境配置。
	Cfg  string
	Ctime int64
	Utime int64
	NextTime int64 `gorm:"index"`

}


type JobDAO interface {

	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context,job domain.Job) error
	UpdateTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

//司空见惯的  实现接口写db实例
type  JobGORMDAO struct {
	 db *gorm.DB
}

func NewJobGORMDAO(db *gorm.DB) *JobGORMDAO {
	return &JobGORMDAO{db: db}
}

const (
	//定义job status
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

func (dao *JobGORMDAO) UpdateNextTime(ctx context.Context, id int64,t time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ? ",id).Updates(
		map[string]any{
			"utime": now,
			"next_time": t.UnixMilli(),
		}).Error
}

func (dao *JobGORMDAO) UpdateTime(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ? ",id).Updates(
		map[string]any{
			"utime": now,
		}).Error
}
func (dao *JobGORMDAO) Preempt(ctx context.Context) (Job, error) {
	//简化写法
	db := dao.db.WithContext(ctx)

	//乐观锁做法先查询
	// 因为找不到，所以可以一直使用for去做，不断抢占。这个东西有东西啊。
	for  {
		//乐观，直接查询
		// 先定义接受实体。
		var j Job
		now := time.Now().UnixMilli()
		//作业
		//缺少找到续约失败的job出来执行 （status = 1 and utime ?     通过utime判断有没有续约
		err := db.Where("status = ? and next_time < ? ",
				jobStatusWaiting,now).First(&j).Error
		if err != nil {
			return j, err
		}
		//走到这里,job有数据
		res := db.Model(&Job{}).Where("id = ? AND  version = ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		// 查询出错
		if res.Error  != nil {
			return Job{}, res.Error
		}

		if res.RowsAffected == 0 {
			//没抢到
			continue
		}
		return j,err
	}
}

func (dao *JobGORMDAO) Release(ctx context.Context, job domain.Job) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&job).
		Where("id = ? ",job.Id).Updates(
			map[string]any{
				"status": jobStatusPaused,
				"utime": now,
			}).Error
}


