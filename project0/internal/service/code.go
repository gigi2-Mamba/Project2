package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"project0/internal/repository"
	"project0/internal/service/sms"
)

// Send生成一个随机验证码，并发送

const secret = "liuxuejin.com"

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	SendFaker(ctx context.Context, biz, phone string) (error, string)
	//CodeGenerate() string
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
	//CodeCrypto(code string) string

}
type codeService struct {
	repos repository.CodeRepository
	sms   sms.Service
}

func NewCodeService(code repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repos: code,
		sms:   sms,
	}
}

// 这里是实际腾讯云短信的实现,我用了sendFaker先实现，后续实际有企业资质再实现

//	func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
//		code := svc.CodeGenerate()
//		err := svc.repos.Set(ctx, biz, phone, code)
//		if err != nil {
//			return err
//		}
//		const codeTplId = "187756"
//
//		return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
//	}

func (svc *codeService) SendFaker(ctx context.Context, biz, phone string) (error, string) {
	code := svc.CodeGenerate()
	//err := svc.repos.Set(ctx, biz, phone, code)
	//if err != nil {
	//	return err
	//}
	//const codeTplId = "187756"
	code0 := svc.CodeCrypto(code)
	err := svc.repos.Set(ctx, biz, phone, code0)
	if err != nil {
		log.Println("test --")
		return err, ""
	}
	const codeTplId = "187756"
	// 先不关心err
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {
		log.Println("localsms failed ", err)
		return err, ""
	}
	return nil, code
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	inputCode0 := svc.CodeCrypto(inputCode)
	ok, err := svc.repos.Verify(ctx, biz, phone, inputCode0)

	if err == repository.ErrCodeVerifyTooMany {
		log.Println("service/code.go ,Code verify too many")
		return false, err
	}
	return ok, err
}

func (svc *codeService) CodeCrypto(code string) string {
	// we need to use some algorithm to crypto verify code
	// the storage package
	engine := md5.New()
	// 实现了一个写方法接口。  对字节内容写入，是否产生错误和写入了多少个字节。 通顶接口
	engine.Write([]byte(secret))
	return hex.EncodeToString([]byte(code))

}
func (svc *codeService) CodeGenerate() string {
	code := rand.Intn(1000000)
	// 06就是6个前导0
	return fmt.Sprintf("%06d", code)
}
