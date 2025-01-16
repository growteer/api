package authn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/growteer/api/infrastructure/solana"
	"github.com/growteer/api/pkg/web3util"
)

const nonce_length = 32

func (s *Service) Login(ctx context.Context, did *web3util.DID, message string, signature string) (sessionToken string, refreshToken string, err error) {
	err = s.verifySignature(ctx, did, message, signature)
	if err != nil {
		return "", "", fmt.Errorf("verification of the received signature failed: %w", err)
	}

	_, err = s.userRepo.GetByDID(ctx, did)
	if err != nil {
		return "", "", fmt.Errorf("could not find an existing user to authenticate: %w", err)
	}

	return s.createNewTokens(ctx, did)
}

func (s *Service) verifySignature(ctx context.Context, did *web3util.DID, message string, signature string) error {
	nonce, err := s.authRepo.GetNonceByDID(ctx, did)
	if err != nil {
		return err
	}

	if !strings.Contains(message, nonce) {
		return fmt.Errorf("message does not contain the correct nonce")
	}

	if err = solana.VerifySignature(message, signature, did.Address); err != nil {
		return err
	}

	return  nil
}

func (s *Service) GenerateNonce(ctx context.Context, did *web3util.DID) (string, error) {
	bytes := make([]byte, nonce_length)

	_, err := rand.Read(bytes)
	if err != nil {
			return "", err
	}

	encoded := hex.EncodeToString(bytes)
	nonce := encoded + ":" + did.Address

	if err = s.authRepo.SaveNonce(ctx, did, nonce); err != nil {
		return "", err
	}

	return nonce, nil
}