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
	topN *atomicx.Value[[]domain.Article]
	ddl *atomicx.Value[time.Time]
	expiration  time.Duration
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	ddl := r.ddl.Load()
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存失效")
	}

	return arts,nil
}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.topN.Store(arts)
	r.ddl.Store(time.Now().Add(r.expiration))
	return nil

}
//搞一个兜底的forceget,当redis崩溃，再次从本地获取
func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	if len(arts) == 0 {
		return []domain.Article{},errors.New("本地缓存失效")
	}
	return arts,nil
}





