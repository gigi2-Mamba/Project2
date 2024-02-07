package sms

import "context"

// 发送短信的抽象
// 屏蔽不同供应商的区别
type Service interface {
	// 上下文暂时又不知道   模版id，短信的内容，短信的接受者
	// 讲到容错机制的时候 context就是为了返回一个超时错误
    // 模版id,numbers要发送的手机号。 args要发送的内容
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
