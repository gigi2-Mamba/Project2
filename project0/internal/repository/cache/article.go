package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project0/internal/domain"
	"time"
)

// Created by Changer on 2024/2/7.
// Copyright 2024 programmer.


type  ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article,error)
	SetFirstPage(ctx context.Context, uid int64,arts []domain.Article) (error)
	DeleteFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article,error)
	Set(ctx context.Context, res domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article,error)
	SetPub(ctx context.Context, res domain.Article) error

}


type  ArticleRedisCache struct {
	 client redis.Cmdable
}

func (r *ArticleRedisCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	key := r.PubKey(id)
	val,err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	//最后一个err直接返回
	err = json.Unmarshal(val,&res)

	return res,err
}

func (r *ArticleRedisCache) SetPub(ctx context.Context, art domain.Article) error {
	key := r.PubKey(art.Id)
	val,err := json.Marshal(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx,key,val,time.Minute * 10).Err()
}

// 这里的get是创作者get自己的创作的帖子
func (r *ArticleRedisCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	key := r.Key(id)

	val,err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	if err != nil {
		return domain.Article{}, err
	}
	return  res,err
}

func (r *ArticleRedisCache) DeleteFirstPage(ctx context.Context, uid int64) error {
      key := r.FirstPageKey(uid)
	return r.client.Del(ctx, key).Err()

}
func (r *ArticleRedisCache) Set(ctx context.Context, res domain.Article) error  {
	key := r.Key(res.Id)
	val,err := json.Marshal(res)
	if err != nil {
		return err
	}
	return r.client.Set(ctx,key,val,time.Second * 10).Err()
}
// 扯到redis 永远和key挂钩，定义key方法
func (r *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	//res, err := r.client.Get(ctx, r.FirstPageKey(uid)).Result()  //常规做法
	//这种方法更简单
	val, err := r.client.Get(ctx, r.FirstPageKey(uid)).Bytes()

	if err != nil {
		return nil,err
	}
	var res  []domain.Article
	//最后一个error直接返回
	err = json.Unmarshal(val, &res)

	return res,nil

}

func (r *ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64,arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		// 用来缓存创作者的帖子列表，缓存摘要就行的
		arts[i].Content = arts[i].Abstract()
	}
	key := r.FirstPageKey(uid)
	val ,err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, val, time.Minute*10).Err()

}

func NewArticleRedisCache(client redis.Cmdable) ArticleCache {
	return &ArticleRedisCache{client: client}
}

func (r *ArticleRedisCache) FirstPageKey(key int64)  string{
	return fmt.Sprintf("article:first_page:%s",key)
}

func (r *ArticleRedisCache) Key(id int64) string{
	return fmt.Sprintf("article:detail:%d",id)
}

func (r *ArticleRedisCache) PubKey(id int64) string{
	return fmt.Sprintf("article:pub:detail:%d",id)
}
