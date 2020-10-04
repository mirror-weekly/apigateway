package main

import (
	"github.com/mirror-media/usersrv"
	"github.com/mirror-media/usersrv/config"
)

func main() {

	cfg := config.Conf{
		Address: "0.0.0.0",
		Port:    80,
	}
	server, err := usersrv.NewServer(cfg)
	if err != nil {
		return
	}

	_ = usersrv.SetRoute(server)

	_ = server.Run()
}
