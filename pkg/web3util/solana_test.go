package web3util_test

import (
	"testing"

	"github.com/growteer/api/pkg/web3util"
	"github.com/stretchr/testify/assert"
)

func Test_VerifySolanaPublicKey(t *testing.T) {
	t.Run("valid public key", func(t *testing.T) {
		// given
		publicKey := "2dqNNQ7jcQpjw15F89KoppKjRCD2U1Ei3eBhUJdYwf7p"

		// when
		err := web3util.VerifySolanaPublicKey(publicKey)

		// then
		assert.NoError(t, err)
	})

	t.Run("too long key", func(t *testing.T) {
		// given
		longPublicKey := "ASXcnp9AJmewtxFvK8iBEwnJe6oJbs348jZqtgdu9dNWtMV8bEPqszg"

		// when
		err := web3util.VerifySolanaPublicKey(longPublicKey)

		// then
		assert.Error(t, err)
	})

	t.Run("non-Base58 key", func(t *testing.T) {
		// given
		unencodedPublicKey := "34b25777224a34b562bdd992db51456603958892e18cfe50cef7d3fcf3b8e8c8"

		// when
		err := web3util.VerifySolanaPublicKey(unencodedPublicKey)

		// then
		assert.Error(t, err)
	})
}
