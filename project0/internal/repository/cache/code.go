package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode           string
	ErrCodeSendTooMany   = errors.New("验证码发送太频繁")
	ErrCodeVerifyTooMany = errors.New("发送太频繁")
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

// 设置验证码
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// cmd with ctx,script, keys[], val
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		log.Println("cache code here is ", err)
		return err
	}

	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		//log.Println("设置验证码成功")
		return nil
	}
}

// func (c *CodeCache) Verify(ctx context.Context, biz, phone, code string) error {
// 怎么采取了bool设置，好处是什么
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// cmd with ctx,script, keys[], val
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		log.Println("Cache use lua with err: ", err)
		return false, err
	}
	// 一个返回结果又多种值对应多种情况时就是用switch
	switch res {
	case -2: // 用户输错了
		return false, nil
	case -1: // 超过验证次数
		return false, ErrCodeVerifyTooMany
	default: // 验证成功
		return true, nil
	}
}

func (c *RedisCodeCache) Key(biz, phone string) string {
	// 构造验证码的key
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)

}
