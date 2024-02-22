package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Created by Changer on 2024/2/9.
// Copyright 2024 programmer.


type Interactive struct {
    Id int64
	//<bizid,biz>
	BizId int64
	Biz string

	ReadCnt int64
	LikeCnt int64
	CollectCnt int64
	Utime  int64
	Ctime  int64
}

type  UserLikeBiz struct {
	Id int64
	Uid int64
	BizId int64
	Biz  string
    Status int
	Utime int64
	Ctime int64
}

type UserCollectBiz struct {
	Id int64
	Uid int64
	BizId int64
	Biz  string
	Cid int64
	favorite string
	Utime int64
	Ctime int64
}
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	IncrLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DecrLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectInfo(ctx context.Context, biz UserCollectBiz) error
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz,error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectBiz,error)
	BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive,error)
}

type InteractiveGORMDAO struct {
     db *gorm.DB
}

func (i *InteractiveGORMDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive

	err := i.db.WithContext(ctx).Where("biz = ? AND biz_id IN ?",biz,ids).First(&res).Error

	return res,err
}

func NewInteractiveGORMDAO(db *gorm.DB) InteractiveDAO{
	return &InteractiveGORMDAO{db: db}
}
// 就是查找收藏记录，根据error给上层liked处理
func (i *InteractiveGORMDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	// gorm的查找基本都是定义载体
	var   res UserLikeBiz
	err := i.db.WithContext(ctx).Where("uid = ? and biz = ? and biz_id = ? and status = ? ",
		uid,biz,id,1).First(&res).Error

	return  res, err
}

func (i *InteractiveGORMDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectBiz, error) {
	var res UserCollectBiz
	err := i.db.WithContext(ctx).Where("uid = ? and biz = ? and biz_id = ?",uid,biz,id).Error
	return res,err
}

func (i *InteractiveGORMDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var res Interactive
	err := i.db.WithContext(ctx).Where("biz = ? and biz_id = ? ",biz,id).First(&res).Error
	return res,err
}



func (i *InteractiveGORMDAO) InsertCollectInfo(ctx context.Context, collect UserCollectBiz) error {
	now := time.Now().UnixMilli()
	collect.Ctime = now
	collect.Utime = now
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&collect).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":now,
				"collect_cnt": gorm.Expr("collect_cnt + 1"),
			}),
		}).Create(&Interactive{
			Biz: collect.Biz,
			Utime: collect.Utime,
			BizId: collect.BizId,
			CollectCnt: 1,
		}).Error
	})
}

func (i *InteractiveGORMDAO) IncrLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	// 怎么还要事务处理？  点赞是一个并发的行为，如何并发修改数据库字段呢？
	// 这里更新两个库所以要事务处理。
	now := time.Now().UnixMilli()
	return  i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLikeBiz{
			Biz:   biz,
			Uid:   uid,
			Status: 1,
			BizId: id,
			Ctime: now,
			Utime: now,
		}).Error
		if err != nil {
			return err
		}
		// 更新阅读数量不应该会出现并发问题吗
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":now,
				"like_cnt": gorm.Expr("like_cnt+1"),
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   id,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (i *InteractiveGORMDAO) DecrLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()

	 return  i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("uid = ? and biz_id = ? and biz = ?",uid,id,biz).
			Updates(map[string]any{
				"utime":now,
				"status": 0,
		}).Error

		if err != nil {
			return  err
		}
		return tx.Model(&Interactive{}).Where("biz = ? and biz_id = ? ",biz,id).
			Updates(map[string]any{
				"utime":now,
				"like_cnt": gorm.Expr("like_cnt - 1"),
		}).Error

	})
}
func (i *InteractiveGORMDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewInteractiveGORMDAO(tx)
		for i := 0; i < len(bizs); i++ {
			err := txDAO.IncrReadCnt(ctx, bizs[i], ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (i *InteractiveGORMDAO) IncrReadCnt(ctx context.Context, biz string, id int64) error {
     now := time.Now().UnixMilli()
	 //log.Println("biz ----",biz)
	 //抽象逻辑  upsert  什么时候使用update
	return i.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("read_cnt+1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   id,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error

}

