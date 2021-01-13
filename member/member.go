package member

import (
	"context"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/v4/auth"
	"github.com/mirror-media/mm-apigateway/config"
	"github.com/pkg/errors"
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
			"firebaseID": firebaseID,
			"action":     "delete",
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
