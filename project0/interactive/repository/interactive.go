package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"project0/interactive/domain"
	"project0/interactive/repository/cache"
	dao "project0/interactive/repository/dao"

	"project0/pkg/loggerDefine"
)

// Created by Changer on 2024/2/9.
// Copyright 2024 programmer.

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error
	DecrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	//BatchIncrReadCnt(ctx context.Context, bizs []string, bids []int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

// 操作缓存和持久
type CacheInteractiveRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     loggerDefine.LoggerV1
}

func (c *CacheInteractiveRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	inters, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}

	return slice.Map(inters, func(idx int, src dao.Interactive) domain.Interactive {
		return c.toDomain(src)
	}), nil
}

func (c *CacheInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, id)
	if err == nil {

		return intr, nil
	}
	//log.Println("互动总数没有缓存")
	ie, err := c.dao.Get(ctx, biz, id)

	if err != nil {
		//c.l.Error("")
		return domain.Interactive{}, err
	}
	// 缓存回写
	if err == nil {
		res := c.toDomain(ie)
		err = c.cache.Set(ctx, biz, id, res)
		if err != nil {
			c.l.Error("回写缓存失败",
				loggerDefine.String("biz", biz),
				loggerDefine.Int64("bizId", id),
				loggerDefine.Error(err))
		}
		return res, nil
	}
	return intr, nil
}

func (c *CacheInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)

	switch err {
	case nil:
		return true, nil
	case ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)

	switch err {
	case nil:
		return true, nil
	case ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func NewCacheInteractiveRepository(cache cache.InteractiveCache, dao dao.InteractiveDAO, l loggerDefine.LoggerV1) InteractiveRepository {
	return &CacheInteractiveRepository{cache: cache, dao: dao, l: l}
}

func (c *CacheInteractiveRepository) AddCollectItem(ctx context.Context, biz string, id, cid, uid int64) error {
	err := c.dao.InsertCollectInfo(ctx, dao.UserCollectBiz{
		Biz: biz,
		Cid: cid,
		Uid: uid,
	})

	if err != nil {
		return err
	}
	return c.cache.IncrCollectCnt(ctx, biz, id)
}

func (c *CacheInteractiveRepository) IncrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.IncrLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}

	return c.cache.IncrLikeCnt(ctx, biz, id, uid)
}

func (c *CacheInteractiveRepository) DecrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DecrLikeInfo(ctx, biz, id, uid)

	if err != nil {
		return err
	}

	return c.cache.DecrLikeCnt(ctx, biz, id, uid)
}

func (c *CacheInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// 统计阅读书应该先走数据库
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 这时候要写缓存了
	return c.cache.IncrReadCntIFPresent(ctx, biz, bizId)

}


// 因为拆分微服务注释掉
//func (c *CacheInteractiveRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
//	// 统计阅读书应该先走数据库
//	err := c.dao.BatchIncrReadCnt(ctx, bizs, bizIds)
//
//	if err != nil {
//		return err
//	}
//	// 这时候要写缓存了
//	// 复用原有单个增加阅读数的方法，redis可以顶得住。  差别不大
//	go func() {
//		for i := 0; i < len(bizs); i++ {
//			er := c.cache.IncrReadCntIFPresent(ctx, bizs[i], bizIds[i])
//			if er != nil {
//				c.l.Error("增加阅读数缓存失败,少一个两个无所谓？", loggerDefine.Error(er))
//			}
//
//		}
//	}()
//
//	return nil
//
//}

func (c *CacheInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}

}
