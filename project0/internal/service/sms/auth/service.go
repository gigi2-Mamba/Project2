package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"project0/internal/service/sms"
)

// 装饰器模式,做安全验证。  对调用sms 服务进行安全装饰限制，限制带有指定token的用户才可以使用
// 安全装饰是为了短信服务平台做限制的
type Service struct {
	svc sms.Service
	// 用来放jwt的key
	key []byte // 干嘛key要用byte 切片

}

func (s Service) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claim SMSClaim

	_, err := jwt.ParseWithClaims(tplToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		return err
	}

	return s.svc.Send(ctx, claim.Tpl, args, numbers...)
}

type SMSClaim struct {
	//穿透嵌入，没有具名
	jwt.RegisteredClaims
	Tpl string
}
