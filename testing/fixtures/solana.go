package fixtures

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func GenerateEd25519KeyPair(t *testing.T) (private ed25519.PrivateKey, publicBase58 string) {
	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	require.NoError(t, err)

	private = ed25519.NewKeyFromSeed(seed)
	public, ok := private.Public().(ed25519.PublicKey)
	require.True(t, ok)

	publicBase58 = base58.Encode(public)

	return
}
