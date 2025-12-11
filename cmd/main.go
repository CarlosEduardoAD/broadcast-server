package main

import (
	"github.com/CarlosEduardoAD/broadcast-server/internal/service"
	"github.com/sevlyar/go-daemon"
)

func main() {
	sprit := service.NewDaemon(&daemon.Context{
		PidFileName: "broadcast-server.pid",
		PidFilePerm: 0644,
		LogFileName: "broadcast-server.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[broadcast-server-daemon]"},
	}, []service.Flag{})

	sprit.Spawn()
}
