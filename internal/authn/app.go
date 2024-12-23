package authn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/growteer/api/infrastructure/ethereum"
)

const nonce_length = 32

type Repository interface {
	Save(ctx context.Context, address, nonce string) error
}

type TokenProvider interface {
	NewSessionToken(address string) (string, error)
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

func (s *Service) Login(ctx context.Context, address, message string, signature []byte) (sessionToken string, err error) {
	addressFromSig, err := ethereum.GetAddressFromSignature(message, signature)
	if err != nil {
		return "", err
	}

	if addressFromSig != address {
		slog.Error(fmt.Sprintf("address %s from signature does not match address %s", addressFromSig, address))
		return "", fmt.Errorf("address and signature do not match")
	}

	sessionToken, err = s.tokenProvider.NewSessionToken(address)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (s *Service) GenerateNonce(ctx context.Context, address string) (string, error) {
	bytes := make([]byte, nonce_length)

	_, err := rand.Read(bytes)
	if err != nil {
			return "", err
	}

	encoded := hex.EncodeToString(bytes)
	nonce := encoded + ":" + address

	if err = s.repo.Save(ctx, address, nonce); err != nil {
		return "", err
	}

	return nonce, nil
}