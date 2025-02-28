package web3util_test

import (
	"testing"

	"github.com/growteer/api/pkg/web3util"
	"github.com/growteer/api/testing/fixtures"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VerifySolanaPublicKey(t *testing.T) {
	t.Run("valid public key", func(t *testing.T) {
		// given
		_, pubKey := fixtures.GenerateEd25519KeyPair(t)

		// when
		err := web3util.VerifySolanaPublicKey(pubKey)

		// then
		assert.NoError(t, err)
	})

	t.Run("non-Base58 public key", func(t *testing.T) {
		// given
		_, pubKey := fixtures.GenerateEd25519KeyPair(t)
		pubKeyDecoded, err := base58.Decode(pubKey)
		require.NoError(t, err)

		// when
		err = web3util.VerifySolanaPublicKey(string(pubKeyDecoded))

		// then
		assert.Error(t, err)
	})
}
