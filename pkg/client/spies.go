package client

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type SpyDialer struct {
	CalledWithURL    string
	CalledWithHeader map[string][]string
}

func (s *SpyDialer) Dial(urlStr string, requestHeader map[string][]string) (*websocket.Conn, *http.Response, error) {

	s.CalledWithURL = urlStr
	s.CalledWithHeader = requestHeader

	// Since this is a spy, we won't actually establish a connection.
	return nil, nil, nil
}
