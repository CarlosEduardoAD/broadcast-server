package client

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Dialer interface {
	Dial(urlStr string, requestHeader map[string][]string) (*websocket.Conn, *http.Response, error)
}
