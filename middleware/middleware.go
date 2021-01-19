package middleware

type CtxKey string

const (
	//CtxGinContexKey is the key of a *gin.Context
	CtxGinContexKey CtxKey = "CtxGinContext"
	//CtxFirebaseClient is the key of a *auth.Client
	CtxFirebaseClient CtxKey = "CtxFirebaseClient"
	//CtxFirebaseDatabaseClient is the key of a *db.Client
	CtxFirebaseDatabaseClient CtxKey = "CtxFirebaseDBClient"
)
const (
	// GCtxTokenKey is the key of a token.Token in *gin.Context
	GCtxTokenKey string = "GCtxToken"
	// GCtxUserIDKey is the key of a string of a User ID in *gin.Context
	GCtxUserIDKey string = "GCtxUserID"
)
