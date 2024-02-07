package failover

import (
	"context"
	"log"
	"time"
)

var ReqChan = make(chan *AsyncSendCodeReq,1000)

func AsyncSendCode(hdl *ResponseTimeFailover)  {

	for  {
		   req := <-ReqChan
		   // 做取余就不会
		   hdl.idx = req.Idx
		log.Println("异步发送就绪",hdl.idx)

		   time.Sleep(time.Second * 3)

		err := hdl.Send(req.Ctx, req.TplId, req.Args, req.Numbers...)
		if err != nil {
			log.Println("异步重试出错：",err)
		}
		log.Println("异步发送成功")

	}
}

type AsyncSendCodeReq struct {
	Ctx context.Context
	TplId string
	Args []string
	Numbers []string
	Idx int32
}
