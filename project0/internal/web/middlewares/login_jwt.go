package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"project0/internal/web/ijwt"
)

type LoginJWTMiddlewareBuilder struct {
	Handler ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(handler ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: handler,
	}
}
/// TBCMISS   miss   field inject

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path ///   ctx.Request.URL.Path  指的是把请求url的路由
		//log.Println(path, " --------------------")
		// 对不需要登录校验的路由放行
		if path == "/users/signup" ||
			path =="/hello" ||
			path =="/setcookie" ||
			path == "/users/login" ||
			path == "/users/login_sms/send/code0" ||
			path == "/users/login_sms0" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" ||
			path== "/oauth2/wechat/setcookie" ||
			path== "/oauth2/wechat/getcookie"{
			return
		}
		//tokenStr, ok := ctx.Get("Bearer")  //  这里的get是对ctx里维护信息传递的那个map string-string
		tokenStr := m.Handler.ExtraToken(ctx)
		//fmt.Println("tsr ", tokenStr)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) { // 参数是个接口就要传指针
			return ijwt.JWTKey, nil
		})
		if err != nil {
			// token 不对，是伪造的
			log.Println("伪造吗 ", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.String(http.StatusOK, "会话过期,请重新登录")
			return
		}
		if !token.Valid || token == nil {
			// token解析出来 ，可能是非法或者过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 为什么设置7天登录的token之后就不要UserAgent的校验了？
		if uc.UserAgent != ctx.GetHeader("User-Agent") {

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

        // 在登录校验里加入了 自动刷新机制

		//expireTime, _ := uc.GetExpirationTime()
		//// 剩余时间小于50s
		//if expireTime.Sub(time.Now()) < time.Second*50 {
		//	uc.ExpiresAt = ijwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString(web.JWTKey)
		//	ctx.Header("x-ijwt-token", tokenStr)
		//	if err != nil {
		//		// 这边不要中断，
		//		log.Println(err)
		//	}
		//}

		// 查看token是否无效
		// 作降级策略当redis没有崩溃时做两个验证   严格判定用户有没有主动退出登录
		err = m.Handler.CheckSession(ctx, uc.Ssid)
		if err != nil  {
			//token 无效 或者redis出错
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 走不严格的这条路
		//cnt, err := m.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid)).Result()
		//if cnt > 0 {
		//	//token 无效 或者redis出错
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		ctx.Set("user", uc)
	}
}
