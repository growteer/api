package authn

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/pkg/web3util"
)

type TokenProvider interface {
	NewRefreshToken(did *web3util.DID) (string, error)
	NewSessionToken(did *web3util.DID) (string, error)
	ParseRefreshToken(token string) (claims *jwt.RegisteredClaims, err error)
	ParseSessionToken(token string) (*jwt.RegisteredClaims, error)
}

type AuthnRepository interface {
	GetNonceByDID(ctx context.Context, did *web3util.DID) (string, error)
	SaveNonce(ctx context.Context, did *web3util.DID, nonce string) error
	GetRefreshTokenByDID(ctx context.Context, did *web3util.DID) (string, error)
	SaveRefreshToken(ctx context.Context, did *web3util.DID, token string) error
}

type ProfileRepository interface {
	Exists(ctx context.Context, did *web3util.DID) bool
}
