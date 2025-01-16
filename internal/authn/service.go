package authn

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/growteer/api/internal/profiles"
	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

type Repository interface {
	GetNonceByDID(ctx context.Context, did *web3util.DID) (string, error)
	SaveNonce(ctx context.Context, did *web3util.DID, nonce string) error
	GetRefreshTokenByDID(ctx context.Context, did *web3util.DID) (string, error)
	SaveRefreshToken(ctx context.Context, did *web3util.DID, token string) error
}

type UserRepository interface {
	GetByDID(ctx context.Context, did *web3util.DID) (*profiles.Profile, error)
}

type TokenProvider interface {
	NewSessionToken(did *web3util.DID) (string, error)
	NewRefreshToken(did *web3util.DID) (string, error)
	ParseRefreshToken(token string) (claims *jwt.RegisteredClaims, err error)
}

type Service struct {
	authRepo Repository
	tokenProvider TokenProvider
	userRepo UserRepository
}

func NewService(authRepo Repository, tokenProvider TokenProvider, userRepo UserRepository) *Service {
	return &Service{
		authRepo: authRepo,
		tokenProvider: tokenProvider,
		userRepo: userRepo,
	}
}

func (s *Service) createNewTokens(ctx context.Context, did *web3util.DID) (newSessionToken string, newRefreshToken string, err error) {
	newSessionToken, err = s.tokenProvider.NewSessionToken(did)
	if err != nil {
		return "", "", gqlutil.InternalError(ctx, "could not generate new session token", err)
	}

	newRefreshToken, err = s.tokenProvider.NewRefreshToken(did)
	if err != nil {
		return "", "", gqlutil.InternalError(ctx, "could not generate new refresh token", err)
	}

	if err = s.authRepo.SaveRefreshToken(ctx, did, newRefreshToken); err != nil {
		return "", "", gqlutil.InternalError(ctx, "could not save new refresh token", err)
	}

	return newSessionToken, newRefreshToken, nil
}
