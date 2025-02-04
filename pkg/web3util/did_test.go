package web3util_test

import (
	"testing"

	"github.com/growteer/api/pkg/web3util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	// given
	address := "randomaddress"
	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, address)

	// when
	didString := did.String()

	// then
	assert.Equal(t, "did:pkh:solana:randomaddress", didString)
}

func Test_DIDFromString(t *testing.T) {
	t.Run("valid pkh solana did", func(t *testing.T) {
		// given
		rawDID := "did:pkh:solana:2dqNNQ7jcQpjw15F89KoppKjRCD2U1Ei3eBhUJdYwf7p"

		// when
		parsed, err := web3util.DIDFromString(rawDID)

		// then
		require.NoError(t, err)
		require.Equal(t, web3util.DIDMethodPKH, parsed.Method)
		require.Equal(t, web3util.NamespaceSolana, parsed.Namespace)
		require.Equal(t, "2dqNNQ7jcQpjw15F89KoppKjRCD2U1Ei3eBhUJdYwf7p", parsed.Address)
	})

	t.Run("address missing", func(t *testing.T) {
		// given
		rawDID := "did:pkh:solana"

		// when
		_, err := web3util.DIDFromString(rawDID)

		// then
		assert.Error(t, err)
	})

	t.Run("unsupported did method", func(t *testing.T) {
		// given
		rawDID := "did:key:solana:2dqNNQ7jcQpjw15F89KoppKjRCD2U1Ei3eBhUJdYwf7p"

		// when
		_, err := web3util.DIDFromString(rawDID)

		// then
		assert.Error(t, err)
	})

	t.Run("unsupported namespace", func(t *testing.T) {
		// given
		rawDID := "did:pkh:eip:2dqNNQ7jcQpjw15F89KoppKjRCD2U1Ei3eBhUJdYwf7p"

		// when
		_, err := web3util.DIDFromString(rawDID)

		// then
		assert.Error(t, err)
	})

	t.Run("invalid address", func(t *testing.T) {
		// given
		rawDID := "did:pkh:solana:ASXcnp9AJmewtxFvK8iBEwnJe6oJbs348jZqtgdu9dNWtMV8bEPqszg"

		// when
		_, err := web3util.DIDFromString(rawDID)

		// then
		assert.Error(t, err)
	})
}
