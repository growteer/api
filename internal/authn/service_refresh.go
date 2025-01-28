package authn

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (newSessionToken string, newRefreshToken string, err error) {
	claims, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", gqlutil.BadInputError(ctx, "could not parse refresh token", gqlutil.ErrCodeInvalidCredentials, err)
	}

	did, err := web3util.DIDFromString(claims.Subject)
	if err != nil {
		slog.Error(err.Error())
		return "", "", gqlutil.InternalError(ctx, "could not parse did from refresh token", err)
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		return "", "", gqlutil.BadInputError(ctx, "invalid solana address in did", gqlutil.ErrCodeInvalidCredentials, err)
	}

	savedToken, err := s.authRepo.GetRefreshTokenByDID(ctx, did)
	if err != nil {
		return "", "", gqlutil.BadInputError(ctx, "invalid refresh token", gqlutil.ErrCodeInvalidCredentials, err)
	}

	if savedToken != refreshToken {
		return "", "", gqlutil.BadInputError(ctx, "invalid refresh token", gqlutil.ErrCodeInvalidCredentials, fmt.Errorf("refresh token does not match the one in the database"))
	}

	return s.createNewTokens(ctx, did)
}

