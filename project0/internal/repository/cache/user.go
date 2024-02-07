package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project0/internal/domain"
	"time"
)

/*
缓存用户
*/
// 利用了第三方包的错误码来映射程序技术栈如redis和mysql的错误
var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.UserProfile, error)
	Set(ctx context.Context, du domain.UserProfile) error
}

// 扩展另一种实现
type MemoryCache struct {
}

// 不用client的做法，面相接口实现，而不是具体实现像建立一个客户端的连接 cli
type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration // 复习下go的time使用  time.duration 指的具体时间单位 秒分xxx,阿拉伯数字是以纳秒为基本单位
}

// 从redis获取profile的缓存
func (c *RedisUserCache) Get(ctx context.Context, id int64) (domain.UserProfile, error) {
	key := c.key(id)
	result, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.UserProfile{}, err
	}
	var u domain.UserProfile
	// 原生json包反序列化，any传指针。 any 本身是一个interface
	err = json.Unmarshal([]byte(result), &u)
	return u, err
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.UserProfile) error {
	key := c.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}

	err = c.cmd.Set(ctx, key, data, c.expiration).Err()

	return err
}

// 收耦合
// 面向接口编程，一定不要去初始化你需要的东西。让外面传进来
func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 70, // 怎么在这里设置了一个过期15分钟
	}
}

// 收耦合
// 面向接口编程，一定不要去初始化你需要的东西。让外面传进来
//func NewOld(addr string) *RedisUserCache {
//	cmd := redis.NewClient(&redis.Options{
//		Addr: addr,
//	})
//
//	return &RedisUserCache{
//		cmd:        cmd,
//		expiration: time.Minute * 15,
//	}
//}
