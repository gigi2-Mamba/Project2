package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"project0/internal/domain"
	"time"
)

/*
User: society-programmer on
Date: 2024/2/22  周四
Time: 11:16
*/

type RankingLocalCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	ddl := r.ddl.Load()
	//首次获取本地缓存，如果本地缓存已经失效了，返回失效错误
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存失效")
	}

	return arts, nil
}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.topN.Store(arts)
	r.ddl.Store(time.Now().Add(r.expiration))
	return nil

}

// 搞一个兜底的forceget,当redis崩溃，再次从本地获取,再次从本地缓存获取就不做失效校验了。
func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	if len(arts) == 0 {
		return []domain.Article{}, errors.New("本地缓存失效")
	}
	return arts, nil
}
