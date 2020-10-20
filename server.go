package usersrv

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/mirror-media/usersrv/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type Server struct {
	conf        *config.Conf
	Engine      *gin.Engine
	FirebaseApp *firebase.App
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func (s *Server) Run() error {
	return s.Engine.Run(fmt.Sprintf("%s:%d", s.conf.Address, s.conf.Port))
}

const firebaseCredentialFile = ""

func NewServer(c config.Conf) (*Server, error) {

	engine := gin.Default()

	opt := option.WithCredentialsFile(firebaseCredentialFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	s := &Server{
		conf:        &c,
		Engine:      engine,
		FirebaseApp: app,
	}
	return s, nil
}
