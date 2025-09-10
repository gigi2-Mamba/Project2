package ijwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
	"time"
)

// how  to design such method?
type RedisJWTHandler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration // rc什么意思
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		client:        client,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 7 * 24}
}

var JWTKey = []byte("oDhIbNhVlYcOtAqNvVaMlFbQrDdObWxT")
var RefreshKey = []byte("oDhIbNhVlYcOtAqNvVaMlFbQrDdObWxg")

// 这个写法更加优雅
var _ Handler = &RedisJWTHandler{}

type UserClaims struct {
	// 将想要的数据放在UserClaims
	// OR standerClaims?
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}

type RefreshCliams struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func (h *RedisJWTHandler) SetLoginJWTToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setRefreshJWT(ctx, uid, ssid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	return h.SetJWT(ctx, uid, ssid)
}

// 退出登录,把长短token都设置非法
func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	// gin context中获取用户
	uc := ctx.MustGet("user").(UserClaims)
	return h.client.Set(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid), "", h.rcExpiration).Err()
}

func (h *RedisJWTHandler) SetJWT(ctx *gin.Context, uid int64, ssid string) error {

	uc := UserClaims{
		Uid: uid,
		// 干货操作方法 思想从http头部获取相应字段的信息     gin.context.GetHeader("key")
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 1 分钟过期  ExpiresAt  jwt在何时过期
			// 不是测试情况设置7*24
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 18)),
		},
		Ssid: ssid,
	}
	// 指定签名方法,和自定义负载，也称用户负载
	token := jwt.NewWithClaims(h.signingMethod, uc) // 这个token不能直接返回，这是一个结构体
	tokenStr, err := token.SignedString(JWTKey)     // 对签名做字符串化
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		log.Println("加密refresh-token出现错误 ", err)
		return err
	}
	// 返回token，通过写入ResponseWriter的头部，也就是响应标头
	ctx.Header("x-jwt-token", tokenStr) //  设置头部，设置网络请求上下文，http 头部携带的字段： “key”-value, value also is string
	// 这里既没有设置ctx.Set() 也没有reids set
	// 可能也不需要
	//err = h.client.Set(ctx, fmt.Sprintf("users:ssid:%s", ssid), uid, h.rcExpiration).Err()
	//if err != nil {
	//	fmt.Println("set jwt token redis fail: ", err.Error())
	//}
	return nil
}

// 这里变小写了 因为只在一个地方用
func (h *RedisJWTHandler) setRefreshJWT(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshCliams{
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(h.signingMethod, &rc)
	tokenstr, err := token.SignedString(RefreshKey)
	if err != nil {
		return err
	}

	ctx.Header("x-refresh-token", tokenstr)

	return nil

}
func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()

	if err != nil {
		return err
	}

	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

// var _ Handler = (*RedisJWTHandler)(nil)
func (rj *RedisJWTHandler) ExtraToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization") //   这个Header 头部 有衣服编码
	if authCode == "" {
		// 没登录
		log.Println("checkLogin JWT version ,没登录")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}

	//segs := strings.Split(authCode, " ") // 编码就是一个空格隔开的吗
	segs := strings.SplitN(authCode, " ", 2)
	if len(segs) != 2 {
		// 没登录 auth 内容是乱传的
		log.Println("checkLogin JWT version ,分段不过关", segs)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	return segs[1]
}
