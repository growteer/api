package authn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (newSessionToken string, newRefreshToken string, err error) {
	claims, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", gqlutil.BadInputError(ctx, "could not parse refresh token", err)
	}

	serializedDid := []byte(claims.Subject)
	var did web3util.DID
	if err = json.Unmarshal(serializedDid, &did); err != nil {
		return "", "", gqlutil.InternalError(ctx, "could not parse did from refresh token", err)
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		return "", "", gqlutil.InternalError(ctx, "invalid solana address in did", err)
	}

	savedToken, err := s.authRepo.GetRefreshTokenByDID(ctx, &did)
	if err != nil {
		return "", "", gqlutil.BadInputError(ctx, "invalid refresh token", err)
	}

	if savedToken != refreshToken {
		return "", "", gqlutil.BadInputError(ctx, "invalid refresh token", fmt.Errorf("refresh token does not match the one in the database"))
	}

	return s.createNewTokens(ctx, &did)
}

