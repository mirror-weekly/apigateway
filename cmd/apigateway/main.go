package main

import (
	"github.com/mirror-media/apigateway"
	"github.com/mirror-media/apigateway/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	// TODO move configs to external config file
	cfg := config.Conf{
		Address:                    "0.0.0.0",
		FirebaseCredentialFilePath: "./firebaseCredential.json",
		Port:                       8080,
		V0RESTfulSrvTargetUrl:      "http://104.199.190.189:8080",
	}
	server, err := apigateway.NewServer(cfg)
	if err != nil {
		return
	}

	err = apigateway.SetRoute(server)
	if err != nil {
		log.Fatalf("error setting up route: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("error runing server: %v", err)
	}
}
