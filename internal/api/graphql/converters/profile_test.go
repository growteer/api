package converters_test

import (
	"context"
	"testing"

	"github.com/growteer/api/internal/api/graphql/converters"
	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/pkg/web3util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	solanaPublicKey    = "3n5zv5r5v5r5v5r5v5r5v5r5v5r5v5r5"
	postalCode         = "12345"
	city               = "New York"
	website            = "https://example.com"
	mailAddress        = "test@example.com"
	validDateOfBirth   = "1990-01-01"
	invalidDateOfBirth = "19901-01-32"
)

func Test_ProfileFromOnboardingInput(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		input := &model.NewProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  validDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
			PostalCode:   &postalCode,
			City:         &city,
			Website:      &website,
		}

		// when
		profile, err := converters.ProfileFromOnboardingInput(ctx, did, input)

		// then
		require.NoError(t, err)
		require.Equal(t, did.String(), profile.DID)
		require.Equal(t, input.Firstname, profile.FirstName)
		require.Equal(t, input.Lastname, profile.LastName)
		require.Equal(t, input.PrimaryEmail, profile.PrimaryEmail)
		require.Equal(t, input.Country, profile.Location.Country)
		require.Equal(t, postalCode, profile.Location.PostalCode)
		require.Equal(t, city, profile.Location.City)
		require.Equal(t, website, profile.Website)
		require.Equal(t, "", profile.About)
		require.Equal(t, "", profile.PersonalGoal)
	})

	t.Run("valid input, only mandatory fields", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		input := &model.NewProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  validDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
		}

		// when
		profile, err := converters.ProfileFromOnboardingInput(ctx, did, input)

		// then
		require.NoError(t, err)
		require.Equal(t, did.String(), profile.DID)
		require.Equal(t, input.Firstname, profile.FirstName)
		require.Equal(t, input.Lastname, profile.LastName)
		require.Equal(t, input.PrimaryEmail, profile.PrimaryEmail)
		require.Equal(t, input.Country, profile.Location.Country)
		require.Equal(t, "", profile.Location.PostalCode)
		require.Equal(t, "", profile.Location.City)
		require.Equal(t, "", profile.Website)
		require.Equal(t, "", profile.About)
		require.Equal(t, "", profile.PersonalGoal)
	})

	t.Run("invalid date of birth", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		input := &model.NewProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  invalidDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
			PostalCode:   &postalCode,
			City:         &city,
			Website:      &website,
		}

		// when
		profile, err := converters.ProfileFromOnboardingInput(ctx, did, input)

		// then
		require.Error(t, err)
		require.Nil(t, profile)

		var badInput apperrors.BadInput
		assert.ErrorAs(t, err, &badInput)
	})
}

func Test_ProfileFromUpdateInput(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		personalGoal := "Personal goal"
		about := "About me"

		input := &model.UpdatedProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  validDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
			PostalCode:   &postalCode,
			City:         &city,
			Website:      &website,
			PersonalGoal: &personalGoal,
			About:        &about,
		}

		// when
		profile, err := converters.ProfileFromUpdateInput(ctx, did, input)

		// then
		require.NoError(t, err)
		require.Equal(t, did.String(), profile.DID)
		require.Equal(t, input.Firstname, profile.FirstName)
		require.Equal(t, input.Lastname, profile.LastName)
		require.Equal(t, input.PrimaryEmail, profile.PrimaryEmail)
		require.Equal(t, input.Country, profile.Location.Country)
		require.Equal(t, postalCode, profile.Location.PostalCode)
		require.Equal(t, city, profile.Location.City)
		require.Equal(t, website, profile.Website)
		require.Equal(t, personalGoal, profile.PersonalGoal)
		require.Equal(t, about, profile.About)
	})

	t.Run("valid input, only mandatory fields", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		input := &model.UpdatedProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  validDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
		}

		// when
		profile, err := converters.ProfileFromUpdateInput(ctx, did, input)

		// then
		require.NoError(t, err)
		require.Equal(t, did.String(), profile.DID)
		require.Equal(t, input.Firstname, profile.FirstName)
		require.Equal(t, input.Lastname, profile.LastName)
		require.Equal(t, input.PrimaryEmail, profile.PrimaryEmail)
		require.Equal(t, input.Country, profile.Location.Country)
		require.Equal(t, "", profile.Location.PostalCode)
		require.Equal(t, "", profile.Location.City)
		require.Equal(t, "", profile.Website)
		require.Equal(t, "", profile.PersonalGoal)
		require.Equal(t, "", profile.About)
	})

	t.Run("invalid date of birth", func(t *testing.T) {
		// given
		ctx := context.Background()
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, solanaPublicKey)
		input := &model.UpdatedProfile{
			Firstname:    "John",
			Lastname:     "Doe",
			DateOfBirth:  invalidDateOfBirth,
			PrimaryEmail: mailAddress,
			Country:      "US",
			PostalCode:   &postalCode,
			City:         &city,
			Website:      &website,
		}

		// when
		profile, err := converters.ProfileFromUpdateInput(ctx, did, input)

		// then
		require.Error(t, err)
		require.Nil(t, profile)

		var badInput apperrors.BadInput
		assert.ErrorAs(t, err, &badInput)
	})
}
