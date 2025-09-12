package websocket

import (
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

// simple forward mechanism
type Hub struct {
	conns *syncx.Map[string, *websocket.Conn] // got the generic Map
}

// 加入连接
func (h *Hub) Add(name string, conn *websocket.Conn) {
	h.conns.Store(name, conn)

	go func() {
		for {
			// 读取数据
			typ, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("err: ", err.Error())
			}
			log.Println("received mesage: ", string(msg))
			// forward
			h.conns.Range(func(key string, value *websocket.Conn) bool {
				if key == name {
					return true
				}
				err := value.WriteMessage(typ, msg)
				if err != nil {
					log.Println("receiver message ,err: ", err.Error())
				}
				return true

			})
		}
	}()
}

// 把一个连接塞入了泛型map里
func TestHub(t *testing.T) {
	upgrader := websocket.Upgrader{}
	// initial hub
	hub := &Hub{conns: &syncx.Map[string, *websocket.Conn]{}}
	//websokcet请求都发到这里
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			writer.Write([]byte("upgrade failed"))
			log.Println("upgrade err: ", err.Error())
			return
		}
		hub.Add(request.URL.Query().Get("name"), conn)

	})

	http.ListenAndServe(":8086", nil)

	select {}
}
