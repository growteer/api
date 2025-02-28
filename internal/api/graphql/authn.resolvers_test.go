package graphql_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/internal/api/graphql"
	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/pkg/web3util"
	"github.com/growteer/api/testing/fixtures"
	"github.com/growteer/api/testing/mocks/internal_/app/authn"
	"github.com/growteer/api/testing/testcontainer"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSessionToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzZXNzaW9uVG9rZW4iLCJpYXQiOjE2MjYwNzQwNzcsImV4cCI6MTYyNjA3NzY3N30.7Q7J9"
	testRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJyZWZyZXNoVG9rZW4iLCJpYXQiOjE2MjYwNzQwNzcsImV4cCI6MTYyNjA3NzY3N30.7Q7J9"
)

func Test_AuthNResolver(t *testing.T) {
	mongoEnv, terminateDB := testcontainer.StartMongoAndGetDetails(t)
	defer terminateDB()

	db := mongodb.NewDB(mongoEnv)
	tokenProvider := authn.NewMockTokenProvider(t)
	resolver := graphql.NewResolver(db, tokenProvider)

	t.Run("Login success, user not onboarded", func(t *testing.T) {
		//given
		privKey, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

		tokenProvider.EXPECT().NewSessionToken(did).Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Return(testRefreshToken, nil)

		//when
		nonceResult, err := resolver.Mutation().GenerateNonce(context.Background(), pubKeyBase58)
		require.NoError(t, err)

		loginDetails := newLoginDetails(privKey, pubKeyBase58, nonceResult.Nonce)
		loginResult, err := resolver.Mutation().Login(context.Background(), loginDetails)
		require.NoError(t, err)

		//then
		assert.False(t, loginResult.IsOnboarded)
		assert.Equal(t, testRefreshToken, loginResult.RefreshToken)
		assert.Equal(t, testSessionToken, loginResult.SessionToken)
	})

	t.Run("Login success, user onboarded", func(t *testing.T) {
		//given
		privKey, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

		tokenProvider.EXPECT().NewSessionToken(did).Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Return(testRefreshToken, nil)

		db.Collection("profiles").InsertOne(context.Background(), map[string]interface{}{
			"_id": web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58).String(),
		})
		defer db.Collection("profiles").DeleteOne(context.Background(), map[string]interface{}{
			"_id": web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58).String(),
		})

		//when
		nonceResult, err := resolver.Mutation().GenerateNonce(context.Background(), pubKeyBase58)
		require.NoError(t, err)

		loginDetails := newLoginDetails(privKey, pubKeyBase58, nonceResult.Nonce)
		loginResult, err := resolver.Mutation().Login(context.Background(), loginDetails)
		require.NoError(t, err)

		//then
		assert.True(t, loginResult.IsOnboarded)
		assert.NotEmpty(t, loginResult.RefreshToken)
		assert.NotEmpty(t, loginResult.SessionToken)
	})

	t.Run("Login fail, invalid address", func(t *testing.T) {
		//given
		privKey, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

		tokenProvider.EXPECT().NewSessionToken(did).Maybe().Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Maybe().Return(testRefreshToken, nil)

		//when
		nonceResult, err := resolver.Mutation().GenerateNonce(context.Background(), pubKeyBase58)
		require.NoError(t, err)

		loginDetails := newLoginDetails(privKey, "invalidAddress", nonceResult.Nonce)
		loginResult, err := resolver.Mutation().Login(context.Background(), loginDetails)

		//then
		require.ErrorAs(t, err, &apperrors.BadInput{})
		assert.Nil(t, loginResult)
	})

	t.Run("Login fail, nonce not found", func(t *testing.T) {
		//given
		privKey, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

		tokenProvider.EXPECT().NewSessionToken(did).Maybe().Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Maybe().Return(testRefreshToken, nil)

		//when
		loginDetails := newLoginDetails(privKey, pubKeyBase58, "invalidNonce")
		loginResult, err := resolver.Mutation().Login(context.Background(), loginDetails)

		//then
		require.ErrorAs(t, err, &apperrors.BadInput{})
		assert.Nil(t, loginResult)
	})

	t.Run("Login fail, invalid signature", func(t *testing.T) {
		//given
		privKey, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)

		tokenProvider.EXPECT().NewSessionToken(did).Maybe().Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Maybe().Return(testRefreshToken, nil)

		//when
		nonceResult, err := resolver.Mutation().GenerateNonce(context.Background(), pubKeyBase58)
		require.NoError(t, err)

		loginDetails := newLoginDetails(privKey, pubKeyBase58, nonceResult.Nonce)
		loginDetails.Signature = "invalidSignature"
		loginResult, err := resolver.Mutation().Login(context.Background(), loginDetails)

		//then
		require.ErrorAs(t, err, &apperrors.BadInput{})
		assert.Nil(t, loginResult)
	})

	t.Run("GenerateNonce fail, invalid address", func(t *testing.T) {
		//given
		_, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		pubKeyRaw, err := base58.Decode(pubKeyBase58)
		require.NoError(t, err)

		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, string(pubKeyRaw))

		tokenProvider.EXPECT().NewSessionToken(did).Maybe().Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Maybe().Return(testRefreshToken, nil)

		//when
		nonceResult, err := resolver.Mutation().GenerateNonce(context.Background(), string(pubKeyRaw))

		//then
		require.ErrorAs(t, err, &apperrors.BadInput{})
		assert.Nil(t, nonceResult)
	})

	t.Run("Refresh success", func(t *testing.T) {
		//given
		_, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)
		initialRefreshToken := "eyKMlcCI6MTYyNjA3NzY3N30.7Q7J9"

		db.Collection("refresh_tokens").InsertOne(context.Background(), map[string]interface{}{
			"_id":   did.String(),
			"token": initialRefreshToken,
		})

		tokenProvider.EXPECT().NewSessionToken(did).Return(testSessionToken, nil)
		tokenProvider.EXPECT().NewRefreshToken(did).Return(testRefreshToken, nil)
		tokenProvider.EXPECT().ParseRefreshToken(initialRefreshToken).Return(&jwt.RegisteredClaims{Subject: did.String()}, nil)

		//when
		refreshResult, err := resolver.Mutation().RefreshSession(context.Background(), model.RefreshInput{
			RefreshToken: initialRefreshToken,
		})
		require.NoError(t, err)

		//then
		assert.Equal(t, testRefreshToken, refreshResult.RefreshToken)
		assert.Equal(t, testSessionToken, refreshResult.SessionToken)
	})

	t.Run("Refresh fail, non-existent refresh token", func(t *testing.T) {
		//given
		_, pubKeyBase58 := fixtures.GenerateEd25519KeyPair(t)
		did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)
		initialRefreshToken := "eyKMlcCI6MTYyFjA3NzY3N30.7Q7J9"

		tokenProvider.EXPECT().ParseRefreshToken(initialRefreshToken).Return(&jwt.RegisteredClaims{Subject: did.String()}, nil)

		//when
		refreshResult, err := resolver.Mutation().RefreshSession(context.Background(), model.RefreshInput{
			RefreshToken: initialRefreshToken,
		})

		//then
		require.ErrorAs(t, err, &apperrors.Unauthenticated{})
		assert.Empty(t, refreshResult)
	})
}

func newLoginDetails(privKey ed25519.PrivateKey, pubKey, nonce string) model.LoginDetails {
	message := "message" + nonce
	signature := ed25519.Sign(privKey, []byte(message))

	return model.LoginDetails{
		Address:   pubKey,
		Message:   message,
		Signature: base64.StdEncoding.EncodeToString(signature),
	}
}
