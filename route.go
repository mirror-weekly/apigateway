package usersrv

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const localfile = "/static"

// Set sets the routing for the gin engine
func SetRoute(server *Server) error {

	// Access Auth service from default app
	defaultClient, err := server.App.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	r := server.Engine
	r.Use(static.Serve("/", static.LocalFile(localfile, false)))

	apiRouter := r.Group("/api")

	apiRouter.GET("/verifyToken", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		idToken := authHeader[len(BEARER_SCHEMA):]

		token, err := defaultClient.VerifyIDToken(c, idToken)
		if err != nil {
			apiLogger.Infof("error verifying ID token: %v", err)
			// apiLogger.Infof("token: %v", idToken)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		apiLogger.Infof("Verified ID token: %v\n", token)
		c.Status(http.StatusOK)
	})

	return nil
}
