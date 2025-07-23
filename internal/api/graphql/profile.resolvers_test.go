package graphql_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/internal/api"
	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/pkg/web3util"
	"github.com/growteer/api/testing/fixtures"
	"github.com/growteer/api/testing/mocks/internal_/app/authn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	sessionToken string = "eyBCDQogICAgImRpZCI6ICJkaWQ6c2FsYW5hOk1vZGVsIiwNCiAgICAicHJpbWFyeUlkIjogImRpZDpzaGFyZTpNb2RlbCIsDQogICAgImVtYWlsIjogInRlc3RAZXhhbXBsZS5jb20iDQp9"

	onboardMutation string = `mutation Onboard($firstName: String!, $lastName: String!, $dateOfBirth: String!, $primaryEmail: String!, $country: String!, $postalCode: String) {
		onboard(profile: { firstName: $firstName, lastName: $lastName, dateOfBirth: $dateOfBirth, primaryEmail: $primaryEmail, country: $country, postalCode: $postalCode }) {
			firstName
			lastName
		}
	}`

	profileQuery string = `query Profile($userDID: String!) {
		profile(userDID: $userDID) {
			firstName
			lastName
			dateOfBirth
			primaryEmail
			location {
				country
				postalCode
				city
			}
			about
			personalGoal
			website
		}
	}`

	updateProfileMutation string = `mutation UpdateProfile($firstName: String!, $lastName: String!, $dateOfBirth: String!, $primaryEmail: String!, $country: String!, $postalCode: String) {
		updateProfile(profile: { firstName: $firstName, lastName: $lastName, dateOfBirth: $dateOfBirth, primaryEmail: $primaryEmail, country: $country, postalCode: $postalCode }) {
			firstName
			lastName
			dateOfBirth
			primaryEmail
			location {
				country
				postalCode
				city
			}
			about
			personalGoal
			website
		}
	}`
)

func Test_Onboard(t *testing.T) {
	db := mongodb.NewDB(config.Mongo)

	t.Run("success", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)

		// when
		var onboardResult struct{ Onboard model.Profile }
		err := gqlClient.Post(
			onboardMutation,
			&onboardResult,
			client.Var("firstName", "John"),
			client.Var("lastName", "Doe"),
			client.Var("dateOfBirth", "1990-01-01"),
			client.Var("primaryEmail", "test@example.com"),
			client.Var("country", "US"),
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, "John", onboardResult.Onboard.FirstName)
		assert.Equal(t, "Doe", onboardResult.Onboard.LastName)
	})

	t.Run("fail, invalid input", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)

		// when
		var onboardResult struct{ Onboard model.Profile }
		err := gqlClient.Post(
			onboardMutation,
			&onboardResult,
			client.Var("firstName", "John"),
			client.Var("lastName", "Doe"),
			client.Var("dateOfBirth", "invaliddate"),
			client.Var("primaryEmail", "test@example.com"),
			client.Var("country", "US"),
		)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BAD_INPUT")
		assert.Empty(t, onboardResult.Onboard)
	})
}

func Test_Profile(t *testing.T) {
	db := mongodb.NewDB(config.Mongo)

	t.Run("success", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)
		_ = onboardTestProfile(t, gqlClient)

		// when
		var profileResult struct{ Profile model.Profile }
		err := gqlClient.Post(profileQuery, &profileResult, client.Var("userDID", did.String()))
		require.NoError(t, err)

		// then
		assert.Equal(t, "John", profileResult.Profile.FirstName)
		assert.Equal(t, "Doe", profileResult.Profile.LastName)
		assert.Contains(t, profileResult.Profile.DateOfBirth, "1990-01-01")
		assert.Equal(t, "test@example.com", profileResult.Profile.PrimaryEmail)
		assert.Equal(t, "US", profileResult.Profile.Location.Country)
		assert.Equal(t, "", *profileResult.Profile.Location.PostalCode)
		assert.Equal(t, "", *profileResult.Profile.Location.City)
		assert.Equal(t, "", *profileResult.Profile.About)
		assert.Equal(t, "", *profileResult.Profile.PersonalGoal)
		assert.Equal(t, "", *profileResult.Profile.Website)
	})

	t.Run("fail, profile doesn't exist", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)
		_ = onboardTestProfile(t, gqlClient)

		_, altPubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		altDID := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, altPubKeyBase58)

		// when
		var profileResult struct{ Profile model.Profile }
		err := gqlClient.Post(
			profileQuery,
			&profileResult,
			client.Var("userDID", altDID.String()),
		)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "NOT_FOUND")
		assert.Empty(t, profileResult.Profile)
	})
}

func Test_UpdateProfile(t *testing.T) {
	db := mongodb.NewDB(config.Mongo)

	t.Run("success", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)
		_ = onboardTestProfile(t, gqlClient)

		// when
		var updateResult struct{ UpdateProfile model.Profile }
		err := gqlClient.Post(
			updateProfileMutation,
			&updateResult,
			client.Var("firstName", "Jane"),
			client.Var("lastName", "Doe"),
			client.Var("dateOfBirth", "1990-01-01"),
			client.Var("primaryEmail", "info@example.com"),
			client.Var("country", "US"),
		)
		require.NoError(t, err)

		// then
		assert.Equal(t, "Jane", updateResult.UpdateProfile.FirstName)
		assert.Equal(t, "Doe", updateResult.UpdateProfile.LastName)
		assert.Contains(t, updateResult.UpdateProfile.DateOfBirth, "1990-01-01")
		assert.Equal(t, "info@example.com", updateResult.UpdateProfile.PrimaryEmail)
		assert.Equal(t, "US", updateResult.UpdateProfile.Location.Country)
		assert.Equal(t, "", *updateResult.UpdateProfile.Location.PostalCode)
		assert.Equal(t, "", *updateResult.UpdateProfile.Location.City)
		assert.Equal(t, "", *updateResult.UpdateProfile.About)
		assert.Equal(t, "", *updateResult.UpdateProfile.PersonalGoal)
		assert.Equal(t, "", *updateResult.UpdateProfile.Website)
	})

	t.Run("fail, profile doesn't exist", func(t *testing.T) {
		// given
		did := fixtures.NewDID(t)
		gqlClient := setupGQLClient(t, did, db)
		_ = onboardTestProfile(t, gqlClient)

		// when
		var updateResult struct{ UpdateProfile model.Profile }
		err := gqlClient.Post(
			updateProfileMutation,
			&updateResult,
			client.Var("firstName", "Jane"),
			client.Var("lastName", "Doe"),
			client.Var("dateOfBirth", ""),
			client.Var("primaryEmail", ""),
			client.Var("country", "US"),
		)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BAD_INPUT")
		assert.Empty(t, updateResult.UpdateProfile)
	})
}

func setupGQLClient(t *testing.T, did *web3util.DID, db *mongo.Database) *client.Client {
	tokenProvider := authn.NewMockTokenProvider(t)
	tokenProvider.EXPECT().ParseSessionToken(mock.Anything).Return(&jwt.RegisteredClaims{Subject: did.String()}, nil)

	gqlServer := api.NewServer(environment.Server{HTTPPort: 8080}, db, tokenProvider)
	return client.New(gqlServer.Router, client.Path("/query"), client.AddHeader("Authorization", "Bearer "+sessionToken))
}

func onboardTestProfile(t *testing.T, gqlClient *client.Client) model.Profile {
	var onboardResult struct{ Onboard model.Profile }
	err := gqlClient.Post(
		onboardMutation,
		&onboardResult,
		client.Var("firstName", "John"),
		client.Var("lastName", "Doe"),
		client.Var("dateOfBirth", "1990-01-01"),
		client.Var("primaryEmail", "test@example.com"),
		client.Var("country", "US"),
	)
	require.NoError(t, err)

	return onboardResult.Onboard
}
