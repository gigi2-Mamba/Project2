package repository

import (
	"context"
	"project0/internal/repository/cache"
)

// 根据链路传递的  封装了上层
var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
var ErrCodeSendTooMany = cache.ErrCodeSendTooMany

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (c *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {

	return c.cache.Set(ctx, biz, phone, code)

}

func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {

	return c.cache.Verify(ctx, biz, phone, code)

}
