package middleware

type CtxKey string

const (
	//CtxGinContexKey is the key of a *gin.Context
	CtxGinContexKey CtxKey = "CtxGinContext"
	//CtxFirebaseClientKey is the key of a *auth.Client
	CtxFirebaseClientKey CtxKey = "CtxFirebaseClient"
	//CtxFirebaseDatabaseClientKey is the key of a *db.Client
	CtxFirebaseDatabaseClientKey CtxKey = "CtxFirebaseDBClient"
)
const (
	// GCtxTokenKey is the key of a token.Token in *gin.Context
	GCtxTokenKey string = "GCtxToken"
	// GCtxUserIDKey is the key of a string of a User ID in *gin.Context
	GCtxUserIDKey string = "GCtxUserID"
)
