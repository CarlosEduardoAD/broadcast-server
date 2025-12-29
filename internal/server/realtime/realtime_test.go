package realtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CarlosEduardoAD/broadcast-server/internal/pub"
	"github.com/CarlosEduardoAD/broadcast-server/internal/sub"
	"github.com/gorilla/websocket"
)

func TestWebsocketRoute_EchoesMessage(t *testing.T) {
	// reset publisher and bypass origin checks for tests
	publisher = pub.NewPublisher([]sub.Subscriber{})
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	srv := httptest.NewServer(http.HandlerFunc(WebsocketRoute))
	defer srv.Close()

	wsURL := strings.Replace(srv.URL, "http", "ws", 1) + "/ws"

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer c.Close()

	payload := `{"data":"hello"}`

	if err := c.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
		t.Fatalf("write message: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	mt, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("read message: %v", err)
	}

	if mt != websocket.TextMessage {
		t.Fatalf("expected message type %d, got %d", websocket.TextMessage, mt)
	}

	if string(msg) != payload {
		t.Fatalf("expected payload %s, got %s", payload, string(msg))
	}
}
