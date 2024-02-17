package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
Created by payden-programmer on 2024/2/16.
*/

type UserReadBiz struct {
	Id int64
	Uid int64
	BizId int64
	Biz  string
	Utime int64
	Ctime int64
}

type HistoryDAO interface {
	InsertReadInfo(ctx context.Context,read UserReadBiz) error

}

type HistoryGORMDAO struct {
	db *gorm.DB
}

func (h *HistoryGORMDAO) InsertReadInfo(ctx context.Context, read UserReadBiz) error {
	read.Ctime = read.Utime
	//return h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先插入阅读记录,    阅读这种反复操作就是upsert语句
		return h.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"utime": read.Utime,
				}),
		}).Create(&read).Error


}

func NewHistoryGORMDAO(db *gorm.DB) HistoryDAO {
	return &HistoryGORMDAO{db: db}
}
