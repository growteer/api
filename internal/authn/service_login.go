package authn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/growteer/api/infrastructure/solana"
)

const nonce_length = 32

func (s *Service) Login(ctx context.Context, address, message string, signature string) (sessionToken string, refreshToken string, err error) {
	err = s.verifySignature(ctx, address, message, signature)
	if err != nil {
		return "", "", fmt.Errorf("verification of the received signature failed: %w", err)
	}

	return s.createNewTokens(ctx, address)
}

func (s *Service) verifySignature(ctx context.Context, address string, message string, signature string) error {
	nonce, err := s.repo.GetNonceByAddress(ctx, address)
	if err != nil {
		return err
	}

	if !strings.Contains(message, nonce) {
		return fmt.Errorf("message does not contain the correct nonce")
	}

	if err = solana.VerifySignature(message, signature, address); err != nil {
		return err
	}

	return  nil
}

func (s *Service) GenerateNonce(ctx context.Context, address string) (string, error) {
	bytes := make([]byte, nonce_length)

	_, err := rand.Read(bytes)
	if err != nil {
			return "", err
	}

	encoded := hex.EncodeToString(bytes)
	nonce := encoded + ":" + address

	if err = s.repo.SaveNonce(ctx, address, nonce); err != nil {
		return "", err
	}

	return nonce, nil
}