package sub

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestUserCreatesSubscriberWithValidData(t *testing.T) {
	s := Subscriber{Name: "alice", Ip: "127.0.0.1"}

	if s.Name == "" {
		t.Fatalf("expected subscriber name to be set")
	}

	assert.Equal(t, "127.0.0.1", s.Ip, "expected IP to match the one set")
}

func TestSubscriberCallWritesMessageSuccessfully(t *testing.T) {
	srvConnCh := make(chan *websocket.Conn, 1)
	done := make(chan struct{})

	upgrader := websocket.Upgrader{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
			return
		}

		srvConnCh <- c

		<-done
		c.Close()
	}))

	defer srv.Close()

	wsURL := "ws" + srv.URL[len("http"):]

	clientConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer clientConn.Close()

	serverConn := <-srvConnCh
	defer serverConn.Close()

	s := Subscriber{Name: "bob"}
	msg := []byte("hello world")

	s.Call(websocket.TextMessage, msg, serverConn)

	mt, received, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("read message error: %v", err)
	}

	assert.Equal(t, websocket.TextMessage, mt, "expected message type to match")
	assert.Equal(t, msg, received, "expected received message to match sent message")

	close(done)
}
