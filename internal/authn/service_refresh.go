package authn

import (
	"context"
	"fmt"
)

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (newSessionToken string, newRefreshToken string, err error) {
	claims, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("could not parse refresh token: %w", err)
	}

	address := claims.Subject
	if !IsValidEthereumAddress(address) {
		return "", "", fmt.Errorf("invalid ethereum address parsed from the refresh token: %s", address)
	}

	savedToken, err := s.repo.GetRefreshTokenByAddress(ctx, address)
	if err != nil {
		return "", "", fmt.Errorf("could not find refresh token for user: %w", err)
	}

	if savedToken != refreshToken {
		return "", "", fmt.Errorf("refresh token doesn't match")
	}

	return s.createNewTokens(ctx, address)
}

