package miniIM

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
define  message type
mini IM,
*/

type WsGateway struct {
	svc *Service
}

func (s *WsGateway) start(addr string) {
	// 依赖一个地址启动
	upgrader := websocket.Upgrader{}

	mux := http.NewServeMux() // 一个加了读写锁的http server, 这就是收获啊  加了读写锁的server为啥这么搞

	mux.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// 升级http请求, 读取内容
		c, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("upgrade error"))
			log.Println("upgrade to websocket fail: ", err.Error())
			return
		}
		//var seq int64
		// where to get uid,jwt / session
		uid := s.Uid(request)
		conn := &Conn{c}
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

	})

	http.ListenAndServe(addr, mux)

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
	Cid int64
}

func (s *WsGateway) Uid(req *http.Request) int64 {
	uidStr := req.Header.Get("uid") //
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid

}
