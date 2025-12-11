package service

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/CarlosEduardoAD/broadcast-server/internal/server"
	"github.com/CarlosEduardoAD/broadcast-server/internal/utils"
	"github.com/CarlosEduardoAD/broadcast-server/pkg/client"
	"github.com/gorilla/websocket"
	"github.com/sevlyar/go-daemon"
)

type Daemon struct {
	Context        *daemon.Context
	Flags          []Flag
	AlreadySpawned bool
}

type Flag struct {
	Name    *string
	Value   string
	Signal  os.Signal
	Handler daemon.SignalHandlerFunc
}

type CustomDialer struct {
	WS *websocket.Dialer
}

func (d *CustomDialer) Dial(urlStr string, requestHeader map[string][]string) (*websocket.Conn, *http.Response, error) {
	return d.WS.Dial(urlStr, requestHeader)
}

func (d *CustomDialer) SetProxy(proxy func(*http.Request) (*url.URL, error)) {
	d.WS.Proxy = proxy
}

var (
	signal = flag.String("s", "", `Send signal to the daemon:
  quit — stop daemon
  stop — stop server
  start — start realtime server`)
	client_signal = flag.String("c", "", `Send signal to the client:
  quit — stop client
  send — send message
  start — start client connection`)

	srv = server.NewServer(8080)
)

func NewDaemon(cntxt *daemon.Context, flags []Flag) *Daemon {
	return &Daemon{
		Context: cntxt,
		Flags:   flags,
	}
}

func (d *Daemon) Spawn() {
	absPath, _ := filepath.Abs("../broadcast-server.pid")
	exists, _ := utils.CheckIfFileExists(absPath)

	if exists {
		log.Println("daemon already spawned")
		return
	}

	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, disconnectHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "start"), syscall.SIGHUP, reloadHandler)
	daemon.AddCommand(daemon.StringFlag(client_signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(client_signal, "send"), syscall.SIGUSR1, disconnectHandler)
	daemon.AddCommand(daemon.StringFlag(client_signal, "connect"), syscall.SIGSTOP, clientConnectHandler)

	if d.AlreadySpawned {
		log.Println("daemon already spawned")
		return
	}

	if len(d.Flags) > 0 {
		d.AddFlags()
	}

	if checkDaemonFlags(d.Context) {
		os.Exit(0)
	}

	dm, err := d.Context.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if dm != nil {
		return
	}
	defer d.Context.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
LOOP:
	for {
		time.Sleep(time.Second)
		select {
		case <-stop:
			srv.Disconnect()
			log.Println("Desconectado")
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

func (d *Daemon) AddFlags() {
	for _, flag := range d.Flags {
		daemon.AddCommand(daemon.StringFlag(flag.Name, flag.Value), flag.Signal, flag.Handler)
	}
}

// checkDaemonFlags retorna true se comandos foram enviados, false caso contrário
func checkDaemonFlags(cntxt *daemon.Context) bool {
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		err = daemon.SendCommands(d)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return true
	}
	return false
}

func clientConnectHandler(sig os.Signal) error {
	log.Println("connecting client...")
	c := client.NewClient(&CustomDialer{WS: websocket.DefaultDialer})
	err := c.Connect("ws://localhost:8080/connect", nil)

	log.Println("funcionou")

	if err != nil {
		log.Println(err)
	}

	return nil
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func disconnectHandler(sig os.Signal) error {
	if sig == syscall.SIGTERM {
		err := srv.Disconnect()
		if err != nil {
			log.Printf("Error disconnecting server: %s", err.Error())
			panic("server disconnection was not successfull!")
		}
	}
	return nil
}

func reloadHandler(sig os.Signal) error {
	srv.Connect()
	return nil
}
