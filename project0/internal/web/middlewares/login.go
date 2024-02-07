package middlewares

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	//gob.  注册类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			log.Println("不需要校验")
			return
		}
		sess := sessions.Default(ctx)
		fmt.Println("ses nil ? ", sess == nil)
		userId := sess.Get("userId")
		fmt.Println("userid isnil", userId)
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Println("here occur with   ERROR! securecookie: error")
			return
		}

		now := time.Now()

		//如何需要刷新
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			sess.Set(updateTimeKey, now) // 下一次set会覆盖这个sess， 所以要重新设置？
			sess.Set("userId", userId)   // 不拿这个会被覆盖吗？
			err := sess.Save()

			if err != nil {
				//打日志
				log.Println("checkLogin err : ", err)
			}
		}

	}
}
