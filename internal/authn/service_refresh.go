package authn

import (
	"context"

	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/pkg/web3util"
)

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (newSessionToken string, newRefreshToken string, err error) {
	claims, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", apperrors.BadInput{
			Code:    apperrors.ErrCodeInvalidInput,
			Message: "could not parse refresh token",
			Wrapped: err,
		}
	}

	did, err := web3util.DIDFromString(claims.Subject)
	if err != nil {
		return "", "", apperrors.Internal{
			Code:    apperrors.ErrCodeInternalError,
			Message: "could not parse did from refresh token",
			Wrapped: err,
		}
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		return "", "", apperrors.Unauthenticated{
			Code:    apperrors.ErrCodeUnauthenticated,
			Message: "invalid solana address in refresh token's did",
			Wrapped: err,
		}
	}

	savedToken, err := s.authRepo.GetRefreshTokenByDID(ctx, did)
	if err != nil {
		return "", "", apperrors.Unauthenticated{
			Code:    apperrors.ErrCodeUnauthenticated,
			Message: "could not get refresh token for the did from database",
			Wrapped: err,
		}
	}

	if savedToken != refreshToken {
		return "", "", apperrors.Unauthenticated{
			Code:    apperrors.ErrCodeUnauthenticated,
			Message: "refresh token does not match the one in the database",
		}
	}

	newSessionToken, newRefreshToken, err = s.createNewTokens(ctx, did)
	if err != nil {
		return "", "", apperrors.Internal{
			Code:    apperrors.ErrCodeInternalError,
			Message: "could not create new tokens",
			Wrapped: err,
		}
	}

	return newSessionToken, newRefreshToken, nil
}
