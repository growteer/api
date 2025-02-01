package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.60

import (
	"context"
	"log/slog"

	"github.com/growteer/api/internal/api/graphql/converters"
	"github.com/growteer/api/internal/api/graphql/gqlutil"
	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/infrastructure/session"
	"github.com/growteer/api/pkg/web3util"
)

// Onboard is the resolver for the onboard field.
func (r *mutationResolver) Onboard(ctx context.Context, profile model.NewProfile) (*model.Profile, error) {
	did, err := session.GetAuthenticatedDID(ctx)
	if err != nil {
		return nil, err
	}

	newProfile, err := converters.ProfileFromSignupInput(ctx, did, &profile)

	savedProfile, err := r.profileService.CreateProfile(ctx, *newProfile)
	if err != nil {
		return nil, gqlutil.InternalError(ctx, err.Error(), err)
	}

	gqlProfileModel := &model.Profile{
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

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, profile model.UpdatedProfile) (*model.Profile, error) {
	did, err := session.GetAuthenticatedDID(ctx)
	if err != nil {
		return nil, err
	}

	profileUpdate, err := converters.ProfileFromUpdateInput(ctx, did, &profile)
	if err != nil {
		return nil, err
	}

	updatedProfile, err := r.profileService.UpdateProfile(ctx, did, profileUpdate)
	if err != nil {
		return nil, err
	}

	return &model.Profile{
		Firstname:    updatedProfile.FirstName,
		Lastname:     updatedProfile.LastName,
		PrimaryEmail: updatedProfile.PrimaryEmail,
		Location: &model.Location{
			Country:    updatedProfile.Location.Country,
			PostalCode: &updatedProfile.Location.PostalCode,
			City:       &updatedProfile.Location.City,
		},
		Website:      &updatedProfile.Website,
		PersonalGoal: &updatedProfile.PersonalGoal,
		About:        &updatedProfile.About,
	}, nil
}

// Profile is the resolver for the profile field.
func (r *queryResolver) Profile(ctx context.Context, userDid string) (*model.Profile, error) {
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

	profileDTO := &model.Profile{
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
