package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

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
	}

	var cfg config.Conf
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	server, err := apigateway.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = apigateway.SetHealthRoute(server)
	if err != nil {
		log.Fatalf("error setting up health route: %v", err)
	}

	err = apigateway.SetRoute(server)
	if err != nil {
		log.Fatalf("error setting up route: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", server.Conf.Address, server.Conf.Port),
		Handler: server.Engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Infof("server listening to %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = errors.Wrap(shutdown(srv), err.Error())
			log.Fatalf("listen: %s\n", err)
		} else if err != nil {
			err = errors.Wrap(shutdown(nil), err.Error())
			log.Fatalf("error server closed: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	if err := shutdown(srv); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}
	os.Exit(0)
}

func shutdown(server *http.Server) error {
	if server != nil {
		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
