// Package member defines the member related functions
package member

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	"github.com/mirror-media/mm-apigateway/graph/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/v4/auth"
	"github.com/mirror-media/mm-apigateway/config"
	"github.com/mirror-media/mm-apigateway/token"
	"github.com/pkg/errors"
)

const (
	MsgAttrKeyAction     = "action"
	MsgAttrKeyFirebaseID = "firebaseID"
)
const (
	MsgAttrValueDelete = "delete"
)

// Delete performs a series of actions to revoke token, remove firebase user and request to disable the member in the DB
func Delete(parent context.Context, c config.Conf, client *auth.Client, firebaseID string) error {
	if err := revokeFirebaseToken(parent, client, firebaseID); err != nil {
		return err
	}
	if err := deleteFirebaseUser(parent, client, firebaseID); err != nil {
		return err
	}
	go publishDeleteMemberMessage(context.Background(), c.ProjectID, c.PubSubTopicIDDeleteMember, firebaseID)
	return nil
}

func revokeFirebaseToken(parent context.Context, client *auth.Client, firebaseID string) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	if err := client.RevokeRefreshTokens(ctx, firebaseID); err != nil {
		log.Errorf("error revoking tokens for user: %v, %v", firebaseID, err)
		return err
	}
	log.Infof("revoked tokens for user: %v", firebaseID)
	// accessing the user's TokenValidAfter
	u, err := client.GetUser(ctx, firebaseID)
	if err != nil {
		log.Errorf("error getting user %s: %v", firebaseID, err)
		return err
	}
	timestamp := u.TokensValidAfterMillis / 1000
	log.Printf("the refresh tokens were revoked at: %d (UTC seconds) ", timestamp)
	return nil
}

func deleteFirebaseUser(parent context.Context, client *auth.Client, firebaseID string) error {

	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	err := client.DeleteUser(ctx, firebaseID)
	if err != nil {
		err = errors.WithMessagef(err, "member(%s) deletion failed", firebaseID)
		log.Error(err)
		return err
	}
	return nil
}

func publishDeleteMemberMessage(parent context.Context, projectID string, topicID string, firebaseID string) error {

	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		err = errors.WithMessage(err, "error creating client for pubsub")
		log.Error(err)
		return err
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			MsgAttrKeyFirebaseID: firebaseID,
			MsgAttrKeyAction:     MsgAttrValueDelete,
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		errors.WithMessage(err, "get published message result has error")
		log.Error(err)
		return err
	}
	log.Printf("Published member deletion message with custom attributes(firebaseID: %s); msg ID: %v", firebaseID, id)
	return nil
}

func SubscribeDeleteMember(parent context.Context, c config.Conf, userSrvToken token.Token) error {
	client, err := pubsub.NewClient(parent, c.ProjectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(c.PubSubSubscribeMember)

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)

	tokenString, err := userSrvToken.GetTokenString()
	if err != nil {
		panic(err)
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: tokenString,
			TokenType:   token.TypeJWT,
		},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	graphqlClient := graphql.NewClient(c.ServiceEndpoints.UserGraphQL, graphql.WithHTTPClient(httpClient))

	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			firebaseID := msg.Attributes[MsgAttrKeyFirebaseID]
			log.Infof("Got message to %s member: %s", msg.Attributes[MsgAttrKeyAction], firebaseID)

			switch msg.Attributes[MsgAttrKeyAction] {
			case MsgAttrValueDelete:
				log.Infof("Request Saleor-mirror to delete member: %s", firebaseID)

				preGQL := []string{"mutation($firebaseId: String!) {", "deleteMember(firebaseId: $firebaseId) {"}

				preGQL = append(preGQL, "success")
				preGQL = append(preGQL, "}", "}")
				gql := strings.Join(preGQL, "\n")

				req := graphql.NewRequest(gql)
				req.Var(MsgAttrKeyFirebaseID, firebaseID)

				// Ask User service to delete the member
				var resp struct {
					DeleteMember *model.DeleteMember `json:"deleteMember"`
				}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err = graphqlClient.Run(ctx, req, &resp); err == nil {
					log.Infof("Successfully delete member(%s)", firebaseID)
					msg.Ack()
				} else {
					log.Errorf("Fail to delete member(%s):%v", firebaseID, err)
				}
				cancel()
			default:
				log.Errorf("action(%s) is not supported", msg.Attributes[MsgAttrKeyAction])
			}
		}
	}()

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()
	// Receive blocks until the context is cancelled or an error occurs.
	log.Infof("Pulling subscription: %s", c.PubSubSubscribeMember)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		err = errors.Wrap(err, "receive failed")
		log.Error(err)
		return err
	}

	return nil

}
