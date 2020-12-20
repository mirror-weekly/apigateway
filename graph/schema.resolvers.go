package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/mirror-media/apigateway/graph/generated"
	"github.com/mirror-media/apigateway/graph/model"
)

func (r *mutationResolver) Profile(ctx context.Context) (*model.ProfileType, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Register(ctx context.Context, email string, username string, password1 string, password2 string) (*model.Register, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyAccount(ctx context.Context, token string) (*model.VerifyAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResendActivationEmail(ctx context.Context, email string) (*model.ResendActivationEmail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendPasswordResetEmail(ctx context.Context, email string) (*model.SendPasswordResetEmail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordReset(ctx context.Context, token string, newPassword1 string, newPassword2 string) (*model.PasswordReset, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordSet(ctx context.Context, token string, newPassword1 string, newPassword2 string) (*model.PasswordSet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PasswordChange(ctx context.Context, oldPassword string, newPassword1 string, newPassword2 string) (*model.PasswordChange, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ArchiveAccount(ctx context.Context, password string) (*model.ArchiveAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteAccount(ctx context.Context, id string) (*model.DeleteUpdate, error) {
	// panic(fmt.Errorf("not implemented"))
	// TODO relay to user service
	gCTX, err := GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if id != gCTX.Value("UserID").(string) {
		return nil, fmt.Errorf("id(%s) is not allowed to be deleted by id(%s)", id, gCTX.Value("UserID"))
	}
	// TODO delete Firebase user
	client, err := FirebaseClientFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = client.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	result := true
	deleteUpdate := &model.DeleteUpdate{
		Success: &result,
		Errors:  nil,
	}
	return deleteUpdate, nil
}

func (r *mutationResolver) UpdateAccount(ctx context.Context, firstName *string, lastName *string) (*model.UpdateAccount, error) {
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
