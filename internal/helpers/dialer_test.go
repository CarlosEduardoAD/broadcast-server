package helpers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDialer_DialSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade error: %v", err)
		}
		defer conn.Close()

		// keep the connection open until client closes
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	url := strings.Replace(srv.URL, "http://", "ws://", 1)
	d := &Dialer{}
	conn, resp, err := d.Dial(url, nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	if conn != nil {
		conn.Close()
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func TestDialer_DialFailure(t *testing.T) {
	// Server that does NOT perform websocket upgrade -> dial should fail
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("no upgrade"))
	}))
	defer srv.Close()

	url := strings.Replace(srv.URL, "http://", "ws://", 1)
	d := &Dialer{}
	conn, resp, err := d.Dial(url, nil)
	assert.Error(t, err)
	assert.Nil(t, conn)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
