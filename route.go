package usersrv

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Set sets the routing for the gin engine
func SetRoute(r *gin.Engine) error {

	r.Use(static.Serve("/", static.LocalFile("/home/baronchiu_mirrormedia_mg/user/static", false)))

	return nil
}
