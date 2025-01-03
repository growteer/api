package authn

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	GetNonceByAddress(ctx context.Context, address string) (string, error)
	SaveNonce(ctx context.Context, address, nonce string) error
	GetRefreshTokenByAddress(ctx context.Context, address string) (string, error)
	SaveRefreshToken(ctx context.Context, address, token string) error
}

type TokenProvider interface {
	NewSessionToken(address string) (string, error)
	NewRefreshToken(address string) (string, error)
	ParseRefreshToken(token string) (claims *jwt.RegisteredClaims, err error)
}

type Service struct {
	repo Repository
	tokenProvider TokenProvider
}

func NewService(repo Repository, tokenProvider TokenProvider) *Service {
	return &Service{
		repo: repo,
		tokenProvider: tokenProvider,
	}
}

func (s *Service) createNewTokens(ctx context.Context, address string) (newSessionToken string, newRefreshToken string, err error) {
	newSessionToken, err = s.tokenProvider.NewSessionToken(address)
	if err != nil {
		return "", "", fmt.Errorf("could not generate new session token: %w", err)
	}

	newRefreshToken, err = s.tokenProvider.NewRefreshToken(address)
	if err != nil {
		return "", "", fmt.Errorf("could not generate new refresh token: %w", err)
	}

	if err = s.repo.SaveRefreshToken(ctx, address, newRefreshToken); err != nil {
		return "", "", fmt.Errorf("failed saving the new refresh token: %w", err)
	}

	return newSessionToken, newRefreshToken, nil
}
