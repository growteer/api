package fixtures

import (
	"testing"

	"github.com/growteer/api/pkg/web3util"
)

func NewDID(t *testing.T) *web3util.DID {
	_, pubKeyBase58 := GenerateEd25519KeyPair(t)
	return web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, pubKeyBase58)
}
