package localsms

import (
	"context"
	"log"
)

type Service2 struct {
}

func NewService2() *Service2 {
	return &Service2{}
}

func (s *Service2) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("2号服务商验证码是", args)
	return nil
}
