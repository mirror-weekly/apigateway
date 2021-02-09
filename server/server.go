// Package server define the necessary component of a server
package server

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
	"github.com/mirror-media/mm-apigateway/config"
	"github.com/mirror-media/mm-apigateway/token"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type ServiceEndpoints struct {
	UserGraphQL string
}

type Server struct {
	Conf                   *config.Conf
	Engine                 *gin.Engine
	FirebaseApp            *firebase.App
	FirebaseClient         *auth.Client
	FirebaseDatabaseClient *db.Client
	Services               *ServiceEndpoints
	UserSrvToken           token.Token
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func (s *Server) Run() error {
	return s.Engine.Run(fmt.Sprintf("%s:%d", s.Conf.Address, s.Conf.Port))
}

func NewServer(c config.Conf) (*Server, error) {

	engine := gin.Default()

	opt := option.WithCredentialsFile(c.FirebaseCredentialFilePath)

	config := &firebase.Config{
		DatabaseURL: c.FirebaseRealtimeDatabaseURL,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing app")
	}

	firebaseClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize thr Firebase Auth Client")
	}

	dbClient, err := app.Database(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize the Firebase Database Client")
	}

	gatewayToken, err := token.NewGatewayToken(c.TokenSecretName, c.ProjectID)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to retrieve the latest token(%s)", c.TokenSecretName)
	}

	s := &Server{
		Conf:                   &c,
		Engine:                 engine,
		FirebaseApp:            app,
		FirebaseClient:         firebaseClient,
		FirebaseDatabaseClient: dbClient,
		Services: &ServiceEndpoints{
			UserGraphQL: c.ServiceEndpoints.UserGraphQL,
		},
		UserSrvToken: gatewayToken,
	}
	return s, nil
}
