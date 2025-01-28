package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.60

import (
	"context"
	"log/slog"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/growteer/api/graph/model"
	"github.com/growteer/api/infrastructure/session"
	"github.com/growteer/api/internal/profiles"
	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

// GenerateNonce is the resolver for the generateNonce field.
func (r *mutationResolver) GenerateNonce(ctx context.Context, input model.NonceInput) (*model.NonceResult, error) {
	if err := web3util.VerifySolanaPublicKey(input.Address); err != nil {
		return nil, gqlutil.BadInputError(ctx, "invalid solana address", gqlutil.ErrCodeInvalidCredentials, err)
	}

	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, input.Address)

	nonce, err := r.authnService.GenerateNonce(ctx, did)
	if err != nil {
		return nil, err
	}

	return &model.NonceResult{
		Nonce: nonce,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResult, error) {
	if err := web3util.VerifySolanaPublicKey(input.Address); err != nil {
		return nil, gqlutil.BadInputError(ctx, "invalid solana address", gqlutil.ErrCodeInvalidCredentials, err)
	}

	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, input.Address)

	sessionToken, refreshToken, err := r.authnService.Login(ctx, did, input.Message, input.Signature)
	if err != nil {
		graphql.AddError(ctx, err)
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh is the resolver for the refresh field.
func (r *mutationResolver) Refresh(ctx context.Context, input *model.RefreshInput) (*model.AuthResult, error) {
	sessionToken, refreshToken, err := r.authnService.RefreshSession(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, input model.SignupInput) (*model.UserProfile, error) {
	did, err := session.GetAuthenticatedDID(ctx)
	if err != nil {
		return nil, err
	}

	dateOfBirth, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		return nil, gqlutil.BadInputError(ctx, "invalidly formatted date of birth", gqlutil.ErrCodeInvalidDateTimeFormat, err)
	}

	location := profiles.Location{
		Country: input.Country,
	}
	if input.PostalCode != nil {
		location.PostalCode = *input.PostalCode
	}
	if input.City != nil {
		location.City = *input.City
	}

	newProfile := profiles.Profile{
		DID:          did.String(),
		FirstName:    input.Firstname,
		LastName:     input.Lastname,
		DateOfBirth:  dateOfBirth,
		PrimaryEmail: input.PrimaryEmail,
		Location:     location,
	}

	if input.Website != nil {
		newProfile.Website = *input.Website
	}

	savedProfile, err := r.profileService.CreateProfile(ctx, newProfile)
	if err != nil {
		return nil, gqlutil.InternalError(ctx, err.Error(), err)
	}

	gqlProfileModel := &model.UserProfile{
		Firstname:    savedProfile.FirstName,
		Lastname:     savedProfile.LastName,
		PrimaryEmail: savedProfile.PrimaryEmail,
		Location: &model.Location{
			Country:    savedProfile.Location.Country,
			PostalCode: &savedProfile.Location.PostalCode,
			City:       &savedProfile.Location.City,
		},
	}

	return gqlProfileModel, nil
}

// UserProfile is the resolver for the userProfile field.
func (r *queryResolver) UserProfile(ctx context.Context, userDid string) (*model.UserProfile, error) {
	did, err := session.GetAuthenticatedDID(ctx)
	if err != nil {
		return nil, err
	}

	parsedUserDID, err := web3util.DIDFromString(userDid)
	if err != nil {
		slog.Warn(err.Error(),
			slog.Attr{Key: "profile", Value: slog.StringValue(userDid)},
			slog.Attr{Key: "user", Value: slog.StringValue(did.String())},
		)

		return nil, gqlutil.BadInputError(ctx, "invalid did provided", gqlutil.ErrCodeInvalidInput, err)
	}

	profile, err := r.profileService.GetProfile(ctx, parsedUserDID)
	if err != nil {
		return nil, err
	}

	profileDTO := &model.UserProfile{
		Firstname:    profile.FirstName,
		Lastname:     profile.LastName,
		PrimaryEmail: profile.PrimaryEmail,
		DateOfBirth:  profile.DateOfBirth.String(),
		Location: &model.Location{
			Country:    profile.Location.Country,
			PostalCode: &profile.Location.PostalCode,
			City:       &profile.Location.City,
		},
		About:        &profile.About,
		PersonalGoal: &profile.PersonalGoal,
		Website:      &profile.Website,
	}

	return profileDTO, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
