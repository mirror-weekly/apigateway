package main

import (
	"github.com/mirror-media/usersrv"
	"github.com/mirror-media/usersrv/config"
	log "github.com/sirupsen/logrus"
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

	err = usersrv.SetRoute(server)
	if err != nil {
		log.Fatalf("error setting up route: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("error runing server: %v", err)
	}
}
