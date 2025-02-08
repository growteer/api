package converters

import (
	"context"
	"time"

	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/growteer/api/internal/entities"
	"github.com/growteer/api/pkg/web3util"
)

func ProfileFromOnboardingInput(ctx context.Context, did *web3util.DID, input *model.NewProfile) (*entities.Profile, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		return nil, apperrors.BadInput{
			Message: "invalid date format",
			Wrapped: err,
		}
	}

	location := entities.Location{
		Country: input.Country,
	}
	if input.PostalCode != nil {
		location.PostalCode = *input.PostalCode
	}
	if input.City != nil {
		location.City = *input.City
	}

	profile := &entities.Profile{
		DID:          did.String(),
		FirstName:    input.Firstname,
		LastName:     input.Lastname,
		DateOfBirth:  dateOfBirth,
		PrimaryEmail: input.PrimaryEmail,
		Location:     location,
	}

	if input.Website != nil {
		profile.Website = *input.Website
	}

	return profile, nil
}

func ProfileFromUpdateInput(ctx context.Context, did *web3util.DID, input *model.UpdatedProfile) (*entities.Profile, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, input.DateOfBirth)
	if err != nil {
		return nil, apperrors.BadInput{
			Message: "invalid date format",
			Wrapped: err,
		}
	}

	location := entities.Location{
		Country: input.Country,
	}
	if input.PostalCode != nil {
		location.PostalCode = *input.PostalCode
	}
	if input.City != nil {
		location.City = *input.City
	}

	profile := &entities.Profile{
		DID:          did.String(),
		FirstName:    input.Firstname,
		LastName:     input.Lastname,
		DateOfBirth:  dateOfBirth,
		PrimaryEmail: input.PrimaryEmail,
		Location:     location,
	}

	if input.Website != nil {
		profile.Website = *input.Website
	}
	if input.PersonalGoal != nil {
		profile.PersonalGoal = *input.PersonalGoal
	}
	if input.About != nil {
		profile.About = *input.About
	}

	return profile, nil
}
