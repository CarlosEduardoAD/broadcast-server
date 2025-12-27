package helpers

import (
	"log"
	"net/http"

	"github.com/CarlosEduardoAD/broadcast-server/pkg/client"
	"github.com/gorilla/websocket"
)

type Dialer struct{}

func (d *Dialer) Dial(urlStr string, requestHeader map[string][]string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(urlStr, requestHeader)
}

func ClientConnect() {
	srvUrl := "ws://localhost:8080/connect"
	clt := client.NewClient(&Dialer{})

	err := clt.Connect(srvUrl, nil)

	if err != nil {
		panic(err)
	}

	for {
		_, message, err := clt.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

}
