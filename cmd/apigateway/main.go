package main

import (
	"os"

	apigateway "github.com/mirror-media/mm-apigateway"
	"github.com/mirror-media/mm-apigateway/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	// name of config file (without extension)
	viper.SetConfigName("config")
	// optionally look for config in the working directory
	viper.AddConfigPath(".")
	// Find and read the config file
	err := viper.ReadInConfig()
	// Handle errors reading the config file
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
		os.Exit(1)
	}

	var cfg config.Conf
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		os.Exit(1)
	}

	server, err := apigateway.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = apigateway.SetHealthRoute(server)
	if err != nil {
		log.Fatalf("error setting up health route: %v", err)
		os.Exit(1)
	}

	err = apigateway.SetRoute(server)
	if err != nil {
		log.Fatalf("error setting up route: %v", err)
		os.Exit(1)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("error runing server: %v", err)
		panic(err)
	}
}
