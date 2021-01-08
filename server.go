package apigateway

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/mirror-media/mm-apigateway/config"
	"github.com/mirror-media/mm-apigateway/token"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type ServiceEndpoints struct {
	UserGraphQL string
}

type Server struct {
	conf           *config.Conf
	Engine         *gin.Engine
	FirebaseApp    *firebase.App
	FirebaseClient *auth.Client
	Services       *ServiceEndpoints
	UserSrvToken   token.Token
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func (s *Server) Run() error {
	return s.Engine.Run(fmt.Sprintf("%s:%d", s.conf.Address, s.conf.Port))
}

func NewServer(c config.Conf) (*Server, error) {

	engine := gin.Default()

	opt := option.WithCredentialsFile(c.FirebaseCredentialFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	firebaseClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("initialization of Firebase Auth Client encountered error: %s", err.Error())
	}

	gatewayToken, err := token.NewGatewayToken(c.TokenSecretName, c.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("retrieving of the latest token(%s) encountered error: %s", c.TokenSecretName, err.Error())
	}

	s := &Server{
		conf:           &c,
		Engine:         engine,
		FirebaseApp:    app,
		FirebaseClient: firebaseClient,
		Services: &ServiceEndpoints{
			UserGraphQL: c.ServiceEndpoints.UserGraphQL,
		},
		UserSrvToken: gatewayToken,
	}
	return s, nil
}
