package authn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/growteer/api/pkg/web3util"
)

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (newSessionToken string, newRefreshToken string, err error) {
	claims, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("could not parse refresh token: %w", err)
	}

	serializedDid := []byte(claims.Subject)
	var did *web3util.DID
	if err = json.Unmarshal(serializedDid, did); err != nil {
		return "", "", fmt.Errorf("could not parse did %s from refresh token: %w", did.String(), err)
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		return "", "", fmt.Errorf("did %s isn't constructed with a valid solana address: %w", did.String(), err)
	}

	savedToken, err := s.authRepo.GetRefreshTokenByDID(ctx, did)
	if err != nil {
		return "", "", fmt.Errorf("could not find refresh token for user: %w", err)
	}

	if savedToken != refreshToken {
		return "", "", fmt.Errorf("refresh token doesn't match")
	}

	return s.createNewTokens(ctx, did)
}

