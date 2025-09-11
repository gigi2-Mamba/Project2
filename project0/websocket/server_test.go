package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	upgrader := websocket.Upgrader{} // 这只是一个媒介,用来装载关于http upgrade 的 parameter

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		//responseHeader 可以不传，为什么？
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("upgrade error"))
			return
		}
		ws := &Ws{conn}
		// You should continuously receive data from  the conn.
		go func() {
			ws.ReadCycle()
		}()

		go func() {
			ticker := time.NewTicker(time.Second)
			for now := range ticker.C {
				ws.Write("来自服务端的数据: " + now.String())
			}
		}()

	})
	http.ListenAndServe(":8086", nil)
}

// 封装会好看点,没有对比啊
type Ws struct {
	conn *websocket.Conn
}

func (w *Ws) Write(msg string) {
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write error: ", err.Error())
	}

}

func (w *Ws) ReadCycle() {
	conn := w.conn
	for {

		typ, msg, err := conn.ReadMessage() // 怎么读取信息，ws是依赖连接对象
		if err != nil {
			// exit loop
			// record log
			return
		} // type有三种类型，TextMessag,BinaryMessage,CloseMessage   pingmessage, pongmessage 对应心跳
		switch typ {
		case websocket.CloseMessage:
			conn.Close()
			return
		case websocket.BinaryMessage, websocket.TextMessage:
			log.Println("read message: ", string(msg))
		default:
			// no operation

		}
	}
}
