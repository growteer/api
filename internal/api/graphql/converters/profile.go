package converters

import (
	"context"
	"time"

	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/internal/profiles"
	"github.com/growteer/api/pkg/web3util"
)

func ProfileFromOnboardingInput(ctx context.Context, did *web3util.DID, input *model.NewProfile) (*profiles.Profile, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		return nil, apperrors.BadInput{
			Message: "invalid date format",
			Wrapped: err,
		}
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

	return &newProfile, nil
}

func ProfileFromUpdateInput(ctx context.Context, did *web3util.DID, input *model.UpdatedProfile) (*profiles.Profile, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		return nil, apperrors.BadInput{
			Message: "invalid date format",
			Wrapped: err,
		}
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
	if input.PersonalGoal != nil {
		newProfile.PersonalGoal = *input.PersonalGoal
	}
	if input.About != nil {
		newProfile.About = *input.About
	}

	return &newProfile, nil
}
