package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	domain "project0/interactive/domain"
	"strconv"
	"time"
)

// Created by Changer on 2024/2/9.
// Copyright 2024 programmer.

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
    ErrKeyNotExist = redis.Nil
)

const (
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
)

type InteractiveCache interface {
	IncrReadCntIFPresent(ctx context.Context, biz string, id int64) error
	IncrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error
	DecrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error
	IncrCollectCnt(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, id int64, res domain.Interactive) error
	//ReadArticleHistory(ctx context.Context, record domain.ReadHistoryRecord) error
}

type interactiveCache struct {
	client redis.Cmdable
}

// 用户个人观看历史记录
//func (i *interactiveCache) ReadArticleHistory(ctx context.Context, record domain.ReadHistoryRecord) error {
//	return nil
//}

func (i *interactiveCache) Set(ctx context.Context, biz string, id int64, res domain.Interactive) error {
	key := i.Key(biz, id)
	log.Println("set interactive cache key: ", key)
	err := i.client.HMSet(ctx, key, fieldReadCnt, res.ReadCnt, fieldLikeCnt,
		res.LikeCnt, fieldCollectCnt, res.CollectCnt).Err()
	if err != nil {
		return err
	}
	//过期时间可以随便设置
	return i.client.Expire(ctx, key, time.Minute*10).Err()
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &interactiveCache{client: client}
}
func (i *interactiveCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := i.Key(biz, id)
	val, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(val) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var res domain.Interactive
	res.ReadCnt, _ = strconv.ParseInt(val["read_cnt"], 10, 64)
	res.CollectCnt, _ = strconv.ParseInt(val["collect_cnt"], 10, 64)
	res.LikeCnt, _ = strconv.ParseInt(val["like_cnt"], 10, 64)

	return res, nil
}

func (i *interactiveCache) IncrCollectCnt(ctx context.Context, biz string, id int64) error {
	key := i.Key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (i *interactiveCache) IncrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error {
	key := i.Key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *interactiveCache) DecrLikeCnt(ctx context.Context, biz string, id int64, uid int64) error {
	key := i.Key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *interactiveCache) IncrReadCntIFPresent(ctx context.Context, biz string, id int64) error {
	key := i.Key(biz, id)
	//复习一下，lua脚本取决于具体逻辑去返回返回值
	//res, err := i.client.Eval(ctx, luaIncrCnt, []string{key}, 1).Int()
	log.Println("异步更新阅读数： ", key)
	//defer
	err := i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
	//log.Println("err  read_cnt : ",err)
	return err
}

func (i *interactiveCache) Key(biz string, id int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, id)
}

func (i *interactiveCache) HistoryKey(biz string, uid int64) string {
	return fmt.Sprintf("%s:history:%d", biz, uid)
}
