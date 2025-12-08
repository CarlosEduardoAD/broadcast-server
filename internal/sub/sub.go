package sub

import (
	"log"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Name string
	Ip   string
	Conn *websocket.Conn
}

func NewSubscriber(name, ip string, conn *websocket.Conn) *Subscriber {
	return &Subscriber{
		Name: name,
		Ip:   ip,
		Conn: conn,
	}
}

func (s *Subscriber) Call(mt int, message []byte, c *websocket.Conn) {
	print(c.RemoteAddr().String())
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Println("write:", err)
	}
}
