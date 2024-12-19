package nonce

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

const nonce_length = 32

type Repository interface {
	Save(ctx context.Context, address, nonce string) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{ repository: repository}
}

func (s *Service) GenerateNonce(ctx context.Context, address string) (string, error) {
	bytes := make([]byte, nonce_length)

	_, err := rand.Read(bytes)
	if err != nil {
			return "", err
	}

	encoded := hex.EncodeToString(bytes)
	nonce := encoded + ":" + address

	if err = s.repository.Save(ctx, address, nonce); err != nil {
		return "", err
	}

	return nonce, nil
}