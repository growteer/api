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
	"github.com/growteer/api/testing/testcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	sessionToken string = "eyBCDQogICAgImRpZCI6ICJkaWQ6c2FsYW5hOk1vZGVsIiwNCiAgICAicHJpbWFyeUlkIjogImRpZDpzaGFyZTpNb2RlbCIsDQogICAgImVtYWlsIjogInRlc3RAZXhhbXBsZS5jb20iDQp9"

	onboardMutation string = `mutation Signup($firstName: String!, $lastName: String!, $dateOfBirth: String!, $primaryEmail: String!, $country: String!) {
		onboard(profile: { firstName: $firstName, lastName: $lastName, dateOfBirth: $dateOfBirth, primaryEmail: $primaryEmail, country: $country }) {
			firstName
			lastName
		}
	}`
)

func Test_Onboard(t *testing.T) {
	mongoEnv, terminateDB := testcontainer.StartMongoAndGetDetails(t)
	defer terminateDB()
	db := mongodb.NewDB(mongoEnv)

	_, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

	tokenProvider := authn.NewMockTokenProvider(t)
	tokenProvider.EXPECT().ParseSessionToken(mock.Anything).Return(&jwt.RegisteredClaims{Subject: did.String()}, nil)

	gqlServer := api.NewServer(environment.ServerEnv{HTTPPort: 8080}, db, tokenProvider)
	gqlClient := client.New(gqlServer.Router, client.Path("/query"), client.AddHeader("Authorization", "Bearer "+sessionToken))

	t.Run("fail, invalid input", func(t *testing.T) {
		//when
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

		//then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BAD_INPUT")
		assert.Empty(t, onboardResult.Onboard)
	})

	t.Run("success", func(t *testing.T) {
		//when
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

		//then
		require.NoError(t, err)
		assert.Equal(t, "John", onboardResult.Onboard.FirstName)
		assert.Equal(t, "Doe", onboardResult.Onboard.LastName)
	})
}
