package authn

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/pkg/web3util"
)

type TokenProvider interface {
	NewRefreshToken(did *web3util.DID) (string, error)
	NewSessionToken(did *web3util.DID) (string, error)
	ParseRefreshToken(token string) (claims *jwt.RegisteredClaims, err error)
	ParseSessionToken(token string) (*jwt.RegisteredClaims, error)
}
