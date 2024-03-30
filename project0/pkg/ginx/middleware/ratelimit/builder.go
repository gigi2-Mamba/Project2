package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project0/pkg/limiter"
)

// 装饰器模式，装饰redis实现的滑动窗口限流器编程一个gin插件
type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

// 这里没有完全写死，有需要可以使用 Prefix方法来注入
func NewBuilder(l limiter.Limiter) *Builder {

	return &Builder{
		//默认就是写死了，可以灵活改变一下
		prefix:  "ip-limiter",
		limiter: l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}
func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			log.Println(err)
			//保守做法，服务治理，引用第三方服务时，第三方不可用，就保守方法，直接限流
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}

}
