package helpers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Dialer struct{}

func (d *Dialer) Dial(urlStr string, requestHeader map[string][]string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(urlStr, requestHeader)
}
