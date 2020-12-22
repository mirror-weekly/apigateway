package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/mirror-media/apigateway/graph/generated"
	"github.com/mirror-media/apigateway/graph/model"
	"github.com/shurcooL/graphql"
)

func (r *mutationResolver) Profile(ctx context.Context) (*model.ProfileType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Member(ctx context.Context) (*model.MemberType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateMember(ctx context.Context, email string, firebaseID string, nickname *string) (*model.CreateMember, error) {
	if _, err := r.IsRequestMatchingRequesterFirebaseID(ctx, firebaseID); err != nil {
		return nil, err
	}

	// Ask User service to delete the member
	gqlClient := graphql.NewClient(r.Resolver.UserSrvURL, nil)

	var mutation struct {
		CreateMember struct {
			Success *graphql.Boolean
		} `graphql:"createMember(email: $email, firebaseId: $firebaseId, nickname: $nickname)"`
	}

	variables := map[string]interface{}{
		"email":      graphql.String(email),
		"firebaseId": graphql.String(firebaseID),
		"nickname":   (*graphql.String)(nickname),
	}

	err := gqlClient.Mutate(context.Background(), &mutation, variables)

	return &model.CreateMember{
		Success: (*bool)(mutation.CreateMember.Success),
	}, err
}

func (r *mutationResolver) UpdateMember(ctx context.Context, address *string, birthday *string, firebaseID string, gender *int, name *string, nickname *string, phone *string, profileImage *string) (*model.UpdateMember, error) {

	if _, err := r.IsRequestMatchingRequesterFirebaseID(ctx, firebaseID); err != nil {
		return nil, err
	}

	// Ask User service to delete the member
	gqlClient := graphql.NewClient(r.Resolver.UserSrvURL, nil)

	var mutation struct {
		UpdateMember struct {
			Success *graphql.Boolean
		} `graphql:"createMember(address: $address, birthday: $birthday, firebaseId: $firebaseId, gender: $gender, name: $name, nickname: $nickname, phone: $phone, profileImage: $profileImage)"`
	}

	// Type conversion
	var gReference *int32
	if gender != nil {
		g := int32(*gender)
		gReference = &g
	}

	variables := map[string]interface{}{
		"address":      (*graphql.String)(address),
		"birthday":     (*graphql.String)(birthday),
		"firebaseID":   (graphql.String)(firebaseID),
		"gender":       (*graphql.Int)((*int32)(gReference)),
		"name":         (*graphql.String)(name),
		"nickname":     (*graphql.String)(nickname),
		"phone":        (*graphql.String)(phone),
		"profileImage": (*graphql.String)(profileImage),
	}

	err := gqlClient.Mutate(context.Background(), &mutation, variables)

	return &model.UpdateMember{
		Success: (*bool)(mutation.UpdateMember.Success),
	}, err
}

func (r *mutationResolver) DeleteMember(ctx context.Context, firebaseID string) (*model.DeleteMember, error) {
	if _, err := r.IsRequestMatchingRequesterFirebaseID(ctx, firebaseID); err != nil {
		return nil, err
	}
	// delete Firebase user
	client, err := FirebaseClientFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = client.DeleteUser(ctx, firebaseID)
	if err != nil {
		return nil, err
	}

	// Ask User service to delete the member
	gqlClient := graphql.NewClient(r.Resolver.UserSrvURL, nil)

	var mutation struct {
		DeleteMember struct {
			Success graphql.Boolean
		} `graphql:"deleteMember(firebaseId: $firebaseId)"`
	}

	variables := map[string]interface{}{
		"firebaseId": graphql.String(firebaseID),
	}
	err = gqlClient.Mutate(context.Background(), &mutation, variables)

	deleteMember := &model.DeleteMember{
		Success: (*bool)(&mutation.DeleteMember.Success),
	}
	return deleteMember, err
}

func (r *mutationResolver) VerifyMember(ctx context.Context, token string) (*model.VerifyAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ArchiveAccount(ctx context.Context, password string) (*model.ArchiveAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendSecondaryEmailActivation(ctx context.Context, email string, password string) (*model.SendSecondaryEmailActivation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifySecondaryEmail(ctx context.Context, token string) (*model.VerifySecondaryEmail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SwapEmails(ctx context.Context, password string) (*model.SwapEmails, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) TokenAuth(ctx context.Context, password string, email *string, username *string) (*model.ObtainJSONWebToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyToken(ctx context.Context, token string) (*model.VerifyToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RevokeToken(ctx context.Context, refreshToken string) (*model.RevokeToken, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllProfile(ctx context.Context) ([]*model.ProfileType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.UserNode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.UserNode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, before *string, after *string, first *int, last *int, email *string, username *string, usernameIcontains *string, usernameIstartswith *string, isActive *bool, statusArchived *bool, statusVerified *bool, statusSecondaryEmail *string) (*model.UserNodeConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
