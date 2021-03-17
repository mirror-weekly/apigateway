package token

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"firebase.google.com/go/v4/auth"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type FirebaseToken struct {
	tokenString    *string
	tokenState     firebaseTokenState
	firebaseClient *auth.Client
}

type firebaseTokenState struct {
	sync.Mutex
	state *string
}

func (ftt *firebaseTokenState) setState(state string) {
	ftt.state = &state
}

func (ft *FirebaseToken) GetTokenString() (string, error) {
	if ft.tokenString == nil {
		return "", errors.New("token is nil")
	}
	return *ft.tokenString, nil
}

func (ft *FirebaseToken) ExecuteTokenStateUpdate() error {
	if ft.tokenString == nil {
		return errors.New("token is nil")
	}
	log.Debugf("ExecuteTokenStateUpdate...(token:%s)", *ft.tokenString)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	log.Debugf("VerifyIDTokenAndCheckRevoked...(token:%s)", *ft.tokenString)
	defer cancel()
	_, err := ft.firebaseClient.VerifyIDTokenAndCheckRevoked(ctx, *ft.tokenString)
	log.Debugf("VerifyIDTokenAndCheckRevoked Result...(token:%s)(err:%v)", *ft.tokenString, err)
	if err != nil {
		ft.tokenState.setState(err.Error())
		return err
	}
	ft.tokenState.setState(OK)
	return nil
}

// GetTokenState will automatically update state if cached state is nil
func (ft *FirebaseToken) GetTokenState() string {
	ft.tokenState.Lock()
	defer ft.tokenState.Unlock()
	if ft.tokenState.state == nil {
		if err := ft.ExecuteTokenStateUpdate(); err != nil {
			log.Info(err)
		}
	}
	return *ft.tokenState.state
}

// NewFirebaseToken creates a token and excute the token state update procedure
func NewFirebaseToken(authHeader string, client *auth.Client, uri string) (Token, error) {

	logger := log.WithFields(log.Fields{
		"uri": uri,
	})
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	const BearerSchema = "Bearer "
	var state, tokenString *string
	logger.Debugf("NewFirebaseToken...(authHeader:%s)", authHeader)
	if authHeader == "" {
		s := "authorization header is not provided"
		state = &s
		logger.Debugf("state is:%s)", state)
	} else if !strings.HasPrefix(authHeader, BearerSchema) {
		s := "Not a Bearer token"
		state = &s
		logger.Debugf("state is:%s)", state)
	} else {
		t := (authHeader)[len(BearerSchema):]
		logger.Debugf("trimming header :%s)", t)
		tokenString = &t
	}
	logger.Debugf("final tokenString...(tokenString:%s)", tokenString)
	firebaseToken := &FirebaseToken{
		firebaseClient: client,
		tokenString:    tokenString,
		tokenState: firebaseTokenState{
			state: state,
		},
	}
	firebaseToken.tokenState.Lock()
	go func() {
		defer firebaseToken.tokenState.Unlock()
		firebaseToken.ExecuteTokenStateUpdate()
	}()
	return firebaseToken, nil
}
