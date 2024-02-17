package repository

import (
	"context"
	"project0/internal/domain"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"time"
)

/*
Created by payden-programmer on 2024/2/16.
*/


type ReadHistoryRepository interface {
	AddRecord(ctx context.Context,record domain.ReadHistoryRecord) error
}



type CacheArticleHistoryRepository struct {
	dao dao.HistoryDAO
	cache cache.InteractiveCache
}

func NewCacheArticleHistoryRepository(dao dao.HistoryDAO, cache cache.InteractiveCache) ReadHistoryRepository {
	return &CacheArticleHistoryRepository{dao: dao, cache: cache}
}


func (c *CacheArticleHistoryRepository) AddRecord(ctx context.Context, record domain.ReadHistoryRecord) error{
	now := time.Now().UnixMilli()
	record.Utime = now

	err := c.dao.InsertReadInfo(ctx, dao.UserReadBiz{
		Uid:   record.Uid,
		Biz:   record.Biz,
		BizId: record.BizId,
		Utime: record.Utime,

	})
	if err != nil {
		return err
	}
	// tbc 用户历史浏览记录
	return  c.cache.ReadArticleHistory(ctx,record)
}
