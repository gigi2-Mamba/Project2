package repository

import (
	"context"
	"project0/internal/domain"
	"project0/internal/repository/cache"
)




type RankingRepository interface {
	ReplaceTopN(ctx context.Context,arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article,error)
}


type CacheRankingRepository struct {
	cache cache.RankingCache
	//v1
	redisCache cache.RankingRedisCache
	localCache cache.RankingLocalCache
}

func NewCacheRankingRepositoryV1(redisCache cache.RankingRedisCache, localCache cache.RankingLocalCache) *CacheRankingRepository {
	return &CacheRankingRepository{redisCache: redisCache, localCache: localCache}
}

func NewCacheRankingRepository(cache cache.RankingCache) *CacheRankingRepository {
	return &CacheRankingRepository{cache: cache}
}


func (c *CacheRankingRepository) GetTopNV1(ctx context.Context) ([]domain.Article,error)  {
   // 先获取本地缓存
	res,err := c.localCache.Get(ctx)
	if err == nil {
		return res,nil
	}
	res, err = c.redisCache.Get(ctx)
	if err != nil {
		return c.localCache.ForceGet(ctx)
	}
	_ =c.cache.Set(ctx,res)
	return res,nil
}
func (c *CacheRankingRepository) GetTopN(ctx context.Context) ([]domain.Article,error)  {

	return c.cache.Get(ctx)
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return c.cache.Set(ctx,arts)
}

// repository 层面操作两个cache
//什么情况下会导致repository复杂
func (c *CacheRankingRepository) ReplaceTopNV1(ctx context.Context, arts []domain.Article) error {
	_=c.localCache.Set(ctx,arts)

	return c.cache.Set(ctx,arts)
}

