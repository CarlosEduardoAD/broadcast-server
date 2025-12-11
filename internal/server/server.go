package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CarlosEduardoAD/broadcast-server/internal/server/realtime"
	"github.com/CarlosEduardoAD/broadcast-server/internal/server/routes/healthcheck"
)

type Server struct {
	s http.Server
}

func NewServer(port int) *Server {
	return &Server{
		s: http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
	}
}

func (srv *Server) Connect() {
	go func() {

		log.Println("server running")

		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthcheck.Healthcheck)
		mux.HandleFunc("/connect", realtime.WebsocketRoute)
		srv.s.Handler = mux
		err := srv.s.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

func (srv *Server) Disconnect() error {
	log.Println("Disconnecting Server...")
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	return nil
}
