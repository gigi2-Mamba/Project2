package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
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
		// You should continuously receive data from  the conn.
		go func() {
			for {
				typ, msg, err := conn.ReadMessage()
				if err != nil {
					return
				} // type有三种类型，TextMessag,BinaryMessage,CloseMessage   pingmessage, pongmessage 对应心跳

			}
		}()
	})
}

// 封装会好看点
type Ws struct {
	conn *websocket.Conn
}
