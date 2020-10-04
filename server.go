package usersrv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mirror-media/usersrv/config"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	conf   *config.Conf
	Engine *gin.Engine
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

	s := &Server{
		conf:   &c,
		Engine: engine,
	}
	return s, nil
}
