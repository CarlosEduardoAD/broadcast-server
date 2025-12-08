package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/CarlosEduardoAD/broadcast-server/internal/pub"
	"github.com/CarlosEduardoAD/broadcast-server/internal/sub"
	"github.com/CarlosEduardoAD/broadcast-server/internal/utils"
	websocket "github.com/gorilla/websocket"
)

var (
	publisher = pub.Pub{
		Subscribers: make([]sub.Subscriber, 0),
	}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: utils.BypassCheck,
}

func echo(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	host, _, _ := net.SplitHostPort(c.RemoteAddr().String())

	sub := sub.Subscriber{
		Name: "RANDOM",
		Ip:   host,
		Conn: c,
	}

	close := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		publisher.Remove(sub)

		err := close(code, text)

		return err
	})

	publisher.Subscribe(sub)

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			log.Println("read:", err)

			if strings.Contains(err.Error(), "abnormal closure") {
				c.Close()

				publisher.Remove(sub)
				break
			} else {
				break
			}
		}

		publisher.Publish(message, mt)
	}
}

func Send_message(w http.ResponseWriter, r *http.Request) {

	log.Printf("qtd de conns: %d", len(publisher.Subscribers))
	for _, v := range publisher.Subscribers {
		log.Printf("Ip conectado: %s", v.Ip)
	}

}

func main() {

	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/health", Send_message)
	http.HandleFunc("/echo", echo)

	go func() {
		log.Println("server running")
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	for _, v := range publisher.Subscribers {
		v.Conn.Close()
		publisher.Remove(v)
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
}
