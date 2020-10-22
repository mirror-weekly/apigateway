package usersrv

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const localfile = "/static"

const (
	UserStateActivated              = 200
	UserStateDisabled               = 100
	UserStateMissingInProvider      = 300
	UserStateMissingInMirrorMedia   = 401
	UserStateRegistrationIncomplete = 402
)

// SetRoute sets the routing for the gin engine
func SetRoute(server *Server) error {

	// Access Auth service from default app
	firebaseClient, err := server.FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}

	r := server.Engine
	r.Use(static.Serve("/", static.LocalFile(localfile, false)))

	apiRouter := r.Group("/api")

	apiRouter.GET("/verifyToken", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		idToken := authHeader[len(BearerSchema):]

		token, err := firebaseClient.VerifyIDToken(c, idToken)
		if err != nil {
			apiLogger.Infof("error verifying ID token: %v", err)
			// apiLogger.Infof("token: %v", idToken)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		apiLogger.Infof("Verified ID token: %v", token)
		c.Status(http.StatusOK)
	})

	apiRouter.GET("/users/:userID/attributes/state", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		firebaseID := c.Param("userID")

		// Get user info from firebase
		firebaseUser, err := firebaseClient.GetUser(c, firebaseID)
		if err != nil {
			apiLogger.Infof("firebase get user(%s) error: %v", firebaseID, err)
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]int{"state": UserStateMissingInProvider})
			return
		}

		// Get user info from db

		type User struct {
			ID    *int
			State *int
			Email *string
		}
		// TODO move db to server
		db, err := NewDB()
		if err != nil {
			apiLogger.Infof("db open error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var user User
		db.Where("firebase_id = ?", firebaseID).First(&user)

		var stateToReturn int
		if firebaseUser.Disabled {
			stateToReturn = UserStateDisabled
		} else if user.ID == nil {
			stateToReturn = UserStateMissingInMirrorMedia
		} else if user.Email == nil || *user.Email == "" {
			stateToReturn = UserStateRegistrationIncomplete
		} else {
			stateToReturn = *user.State
		}
		c.JSON(http.StatusOK, map[string]int{"state": stateToReturn})
	})

	apiRouter.GET("/users/:userID", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		firebaseID := c.Param("userID")

		type User struct {
			ID                    int64
			FirebaseID            string
			Email                 string
			Name                  *string
			Nickname              *string
			Bio                   *string
			State                 int
			Birthday              *time.Time
			ImageID               int64
			Gender                int
			Phone                 *string
			AddressID             int64
			Point                 int
			CreatedAt             time.Time
			UpdatedAt             time.Time
			MembershipValidBefore *time.Time
			MembershipType        int
			MembershipValidAfter  *time.Time
			CreatedByOperator     int64
		}

		db, err := NewDB()
		if err != nil {
			apiLogger.Infof("db open error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var user User
		db.Where("firebase_id = ?", firebaseID).First(&user)
		apiLogger.Infof("firebase_id(%s):%+v", firebaseID, user)
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, user)
	})

	apiRouter.POST("/users", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		type Address struct {
			ID            *int64 `json:"-"`
			Nationality   *string
			State         *string
			City          *string
			ZipCode       *string
			District      *string
			StreetAddress *string
			CreatedAt     *time.Time `json:"-"`
			UpdatedAt     *time.Time `json:"-"`
		}

		type Image struct {
			ID        *int64 `json:"-"`
			URL       *string
			CreatedAt *time.Time `json:"-"`
			UpdatedAt *time.Time `json:"-"`
		}
		type User struct {
			FirebaseID            *string
			Email                 *string
			Name                  *string
			Nickname              *string
			Bio                   *string
			State                 *int
			Birthday              *time.Time
			ImageID               *int64 `json:"-"`
			Image                 *Image
			Gender                *int
			Phone                 *string
			AddressID             *int64 `json:"-"`
			Address               *Address
			Point                 *int
			CreatedAt             *time.Time `json:"-"`
			UpdatedAt             *time.Time `json:"-"`
			MembershipValidBefore *time.Time
			MembershipType        *int
			MembershipValidAfter  *time.Time
			CreatedByOperator     *int64
		}

		var user User

		err = c.BindJSON(&user)
		if err != nil {
			apiLogger.Infof("parsing error: %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// validate user
		if user.FirebaseID == nil {
			apiLogger.Info("firebase_id isn't provided")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else if user.Email == nil {
			apiLogger.Infof("email is not provided for firebase_id(%s)", *user.FirebaseID)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// TODO move to server
		db, err := NewDB()
		if err != nil {
			apiLogger.Infof("db open error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		result := db.Create(&user)

		if result.RowsAffected != 1 {
			apiLogger.Infof("user with firebase_id(%s) is not created: %v", *user.FirebaseID, result.Error)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	})

	apiRouter.DELETE("/users/:userID", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		firebaseID := c.Param("userID")
		_, err := firebaseClient.UpdateUser(context.Background(), firebaseID, (&auth.UserToUpdate{}).Disabled(true))
		if err != nil {
			apiLogger.Infof("Disabling firebase_id(%s) failed: %v", firebaseID, err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// TODO move to server
		db, err := NewDB()
		if err != nil {
			apiLogger.Infof("db open error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		type User struct {
			State int
		}

		db.Model(&User{}).Where("firebase_id = ?", firebaseID).Update("state", UserStateDisabled)

	})

	apiRouter.PATCH("/users/:userID", func(c *gin.Context) {
		apiLogger := log.WithFields(log.Fields{
			"path": c.FullPath(),
		})

		firebaseID := c.Param("userID")

		type User struct {
			Email                 *string
			Name                  *string
			Nickname              *string
			Bio                   *string
			State                 *int
			Birthday              *time.Time
			ImageID               *int64
			Gender                *int
			Phone                 *string
			AddressID             *int64
			Point                 *int
			CreatedAt             *time.Time
			UpdatedAt             *time.Time
			MembershipValidBefore *time.Time
			MembershipType        *int
			MembershipValidAfter  *time.Time
			CreatedByOperator     *int64
		}

		var user User

		err = c.BindJSON(&user)
		if err != nil {
			apiLogger.Infof("user with firebase_id(%s) is not updated: %v", firebaseID, err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Get user info from firebase
		firebaseUser, err := firebaseClient.GetUser(c, firebaseID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else if firebaseUser == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]int{"state": UserStateMissingInProvider})
			return
		}

		now := time.Now()
		user.CreatedAt = &now
		user.UpdatedAt = &now
		firebaseUserToUpdate := &auth.UserToUpdate{}
		if user.Email != nil {
			*user.State = UserStateActivated

			firebaseUserToUpdate.Email(*user.Email)
		} else {
			user.State = nil
		}

		if user.Name != nil {
			firebaseUserToUpdate.DisplayName(*user.Name)
		}

		if user.Phone != nil {
			firebaseUserToUpdate.PhoneNumber(*user.Phone)
		}

		// TODO Update image

		_, err = firebaseClient.UpdateUser(
			context.Background(),
			firebaseID,
			firebaseUserToUpdate,
		)
		if err != nil {
			apiLogger.Errorf("error updating user: %v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// TODO move to server
		db, err := NewDB()
		if err != nil {
			apiLogger.Infof("db open error: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		result := db.Model(&user).Where("firebase_id = ?", firebaseID).Updates(user)

		if result.RowsAffected != 1 {
			apiLogger.Infof("user with firebase_id(%s) is not updated: %v", firebaseID, result.Error)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	})

	return nil
}
