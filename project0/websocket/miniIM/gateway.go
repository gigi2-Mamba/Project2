package miniIM

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"project0/pkg/loggerDefine"
	"project0/pkg/saramax"
	"strconv"
	"time"
)

/*
define  message type
mini IM,
*/

// 这样正太看来就一个网关节点
type WsGateway struct {
	svc        *Service
	client     sarama.Client
	l          loggerDefine.LoggerV1 // 这个logger回头看一下，好像也是适配器模式啥的
	conns      syncx.Map[int64, *Conn]
	instanceID string // 这个group id这么显眼
	// 复用grader
	upgrader *websocket.Upgrader
}

func (s *WsGateway) start(addr string) {
	// 依赖一个地址启动
	//upgrader := websocket.Upgrader{}
	mux := http.NewServeMux() // 一个加了读写锁的http server, 这就是收获啊  加了读写锁的server有什么意图
	mux.HandleFunc("/ws", s.wsHandler)
	err := s.subscribeMsg()
	if err != nil {
		// record log
	}
	http.ListenAndServe(addr, mux)
}

func (s *WsGateway) wsHandler(writer http.ResponseWriter, request *http.Request) {
	// 升级http请求, 读取内容
	c, err := s.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.Write([]byte("upgrade error"))
		log.Println("upgrade to websocket fail: ", err.Error())
		return
	}
	//var seq int64
	// where to get uid,jwt / session
	uid := s.Uid(request)
	conn := &Conn{c}
	s.conns.Store(uid, conn) // notice
	go func() {
		for {
			// 转发到service
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read message fail: ", err.Error())
				return
			}
			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				//消息格式不对
				log.Println("unmarshal message fail: ", err.Error())
				continue
			}
			// send to server
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				// 向服务端发送数据会有高流量所以要goroutine 另开
				err = s.svc.Receive(ctx, uid, msg)
				cancel()
				if err != nil {
					// 后端服务返回失败
					log.Println("receive message fail: ", err.Error())
					conn.WriteMessage(websocket.TextMessage, []byte("handle forward message fail"))

					//continue
				}
			}()
		}
	}()

}

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Message(msg Message) error { // 序列化后写入
	val, _ := json.Marshal(msg)
	return c.Conn.WriteMessage(websocket.TextMessage, val)

}

//前后端交互的数据格式
type Message struct {
	// 前端发送的序列号
	Seq int64 `json:"seq"` // 在webosocket中唯一
	//  标记什么类型的消息
	// 图片，视频
	// {"type": "image",content:""}
	Type string `json:"type"`
	//
	content string `json:"content"`
	// send to who?
	// channel id
	Cid int64 `json:"cid"`
}

// start the consumer
func (s *WsGateway) subscribeMsg() error {
	// NewConsumerGroupFromClient,Client is sarama client
	cg, err := sarama.NewConsumerGroupFromClient(s.instanceID, s.client) // group id  有讲究啊  notice
	if err != nil {
		//record  error
		return err
	}
	go func() { // 这里consumer太潦草，因为没有做什么处理，得在真正调用的时候才体现topic
		err := cg.Consume(context.Background(), []string{eventName}, saramax.NewHandler[Event](s.l, s.Consume))
		if err != nil {
			//记录日志
		}
	}()
	return nil
}

func (s *WsGateway) Uid(req *http.Request) int64 { // 利用请求头来区分不同的用户
	uidStr := req.Header.Get("uid") //
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

// 这里的consume 只是说针对一个event,你怎么操作或者event
func (s *WsGateway) Consume(msg *sarama.ConsumerMessage, evt Event) error {
	// i need to consume
	conn, ok := s.conns.Load(evt.Receiver)
	if !ok { // 针对没有链上的节点做测试
		return nil
	}
	val, _ := json.Marshal(evt.Msg)
	// 就只是纯发
	return conn.WriteMessage(websocket.TextMessage, val) // different from old version conn.send
}
