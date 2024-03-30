package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"log"
	"net/http"
	"project0/internal/service"
	"project0/internal/service/oauth2/wechat"
	"project0/internal/web/ijwt"
)

// 接入微信扫码登录，会抽象什么通用能力出来呢

// 需要注入服务
type OAuth2Handler struct {
	svc             wechat.Service
	userSvc         service.UserService
	jwtHandler      ijwt.Handler
	key             []byte
	stateCookieName string
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

func NewOAuth2Handler(svc wechat.Service, userSvc service.UserService, jwtHandler ijwt.Handler) *OAuth2Handler {
	return &OAuth2Handler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("oDhIbNhVlYcOtAqNvVaMlFbQrDdObWxT"),
		stateCookieName: "jwt-state",
		jwtHandler:      jwtHandler,
	}
}

func (o *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
	// 微信扫码需要跳转到微信的url，从微信回调的url
	g := server.Group("/oauth2/wechat") // 首次使用两层的url作为group url
	//构造跳转到微信的url,也就是请求url
	g.GET("/authurl", o.Auth2Url)
	//不清楚腾讯，也就是第三方会用什么http方法回调，用any直接万无一失
	g.Any("/callback", o.CallBack)
	g.GET("/setcookie", o.setcookie)
	g.GET("/getcookie", o.getcookie)

}
func (o *OAuth2Handler) setcookie(context *gin.Context) {
	context.SetCookie("cookietst", "123", 300, "/oauth2/wechat/getcookie", "", false, true)
	context.String(http.StatusOK, "设置cookie成功？")

}

func (o *OAuth2Handler) getcookie(context *gin.Context) {
	cookie, err := context.Cookie("cookietst")
	if err != nil {
		context.String(http.StatusOK, "获取失败")
		return
	}
	context.String(http.StatusOK, cookie)
	log.Println("获取cookie应该成功了")
}

func (o *OAuth2Handler) Auth2Url(ctx *gin.Context) {
	state := uuid.New()
	url, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造URL失败",
			Code: 5,
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "服务器异常，state",
			Code: 5,
		})

	}
	log.Println("can arrive here?  ")
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
	log.Println("set cookie ")

}

// 微信回调的url
func (o *OAuth2Handler) CallBack(ctx *gin.Context) {
	cookieheader := ctx.GetHeader("Set-Cookie")
	fmt.Println("dgegnegne ", cookieheader)
	log.Println(cookieheader == "", "HHHHHHH")
	log.Println("cookie 是否设置成功")
	err := o.verifyState(ctx)
	if err != nil {
		log.Println("微信回调错误? ", err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		return
	}
	code := ctx.Query("code")
	//state := ctx.Query("state")

	wechatInfo, err := o.svc.Verify(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		log.Println("微信登录创建用户出错", err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	//err = o.setRefreshJWT(ctx,u.Id)
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// o.SetJWT(ctx, u.Id)

	//if err != nil {
	//	log.Println("加密jwt有误")
	//	ctx.JSON(http.StatusOK, Result{
	//		Msg: "加密jwt有误",
	//	})
	//}
	err = o.jwtHandler.SetLoginJWTToken(ctx, u.Id)
	if err != nil {
		log.Println("setLoginJwt failed: ", err)
		ctx.String(http.StatusOK, "系统错误")
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "微信扫码登录成功",
	})
	return
}

func (o *OAuth2Handler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	log.Println("statecookiename ", o.stateCookieName)
	ck, err := ctx.Cookie(o.stateCookieName)
	log.Println("whether is cookie get error : ", err)
	if err != nil {
		log.Println("无法获取cooke state,err", err)

		return fmt.Errorf("无法获取cooke state,err %s ", err)
	}
	var sc StateClaims
	// 这里的claims是interface
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		//解析jwt claim,keyfunc决定返回什么东西和错误，  返回要的数据
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析token失败,%s ", err)
	}
	//有人攻击
	if sc.State != state {
		return fmt.Errorf("state 不匹配")
	}
	return nil
}

func (o *OAuth2Handler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {
		log.Println("token err: ", err)
		//ctx.JSON(http.StatusOK, Result{
		//	Msg:  "服务器异常",
		//	Code: 5,
		//})
		return fmt.Errorf("token加密，cookiestate 错误 %s", err.Error())
	}

	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)

	return nil
}
