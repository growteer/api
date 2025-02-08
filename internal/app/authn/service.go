package authn

import (
	"context"

	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/growteer/api/pkg/web3util"
)

type Service struct {
	authRepo      AuthnRepository
	userRepo      ProfileRepository
	tokenProvider TokenProvider
}

func NewService(authRepo AuthnRepository, tokenProvider TokenProvider, userRepo ProfileRepository) *Service {
	return &Service{
		authRepo:      authRepo,
		tokenProvider: tokenProvider,
		userRepo:      userRepo,
	}
}

func (s *Service) createNewTokens(ctx context.Context, did *web3util.DID) (newSessionToken string, newRefreshToken string, err error) {
	newSessionToken, err = s.tokenProvider.NewSessionToken(did)
	if err != nil {
		return "", "", apperrors.Internal{
			Message: "could not generate new session token",
			Wrapped: err,
		}
	}

	newRefreshToken, err = s.tokenProvider.NewRefreshToken(did)
	if err != nil {
		return "", "", apperrors.Internal{
			Message: "could not generate new refresh token",
			Wrapped: err,
		}
	}

	if err = s.authRepo.SaveRefreshToken(ctx, did, newRefreshToken); err != nil {
		return "", "", apperrors.Internal{
			Message: "could not save refresh token",
			Wrapped: err,
		}
	}

	return newSessionToken, newRefreshToken, nil
}
