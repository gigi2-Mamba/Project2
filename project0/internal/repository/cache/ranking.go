package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"project0/internal/domain"
	"time"
)


type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context ) ( []domain.Article,error)
}
type RankingRedisCache struct {
	client redis.Cmdable
	key string
	expiration time.Duration
}

func (r *RankingRedisCache) Get(ctx context.Context) ( []domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

// 这里上线版本要做成动态配置
//动态配置岂不是要单独写一个InitRedisCache.  找机会看看现存的怎么注入这种动态配置。 TBC
func NewRankingRedisCache(client redis.Cmdable) RankingCache {
	return &RankingRedisCache{client: client,
		expiration: time.Minute * 3,
	    key:"ranking:top_n",
	    }
}

func (r *RankingRedisCache)  Set(ctx context.Context, arts []domain.Article) error {
	//把content做成摘要
	for _,art :=range arts {
		art.Content = art.Abstract()
	}

	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx,r.key,val,r.expiration).Err()

}


