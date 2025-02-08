package authn

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/internal/infrastructure/solana"
	"github.com/growteer/api/pkg/web3util"
)

const nonce_length = 32

func (s *Service) Login(ctx context.Context, did *web3util.DID, message string, signature string) (sessionToken string, refreshToken string, err error) {
	err = s.verifySignature(ctx, did, message, signature)
	if err != nil {
		return "", "", apperrors.BadInput{
			Message: "could not verify signature",
			Wrapped: err,
		}
	}

	sessionToken, refreshToken, err = s.createNewTokens(ctx, did)
	if err != nil {
		return "", "", apperrors.Internal{
			Message: "could not create tokens",
			Wrapped: err,
		}
	}

	if !s.userRepo.Exists(ctx, did) {
		err = apperrors.NotFound{
			Message: "user not sign",
		}
	}

	return sessionToken, refreshToken, err
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

	return nil
}

func (s *Service) GenerateNonce(ctx context.Context, did *web3util.DID) (string, error) {
	bytes := make([]byte, nonce_length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", apperrors.Internal{
			Code:    apperrors.ErrCodeInternalError,
			Message: "could not generate nonce",
			Wrapped: err,
		}
	}

	encoded := hex.EncodeToString(bytes)
	nonce := encoded + ":" + did.Address

	if err = s.authRepo.SaveNonce(ctx, did, nonce); err != nil {
		return "", apperrors.Internal{
			Code:    apperrors.ErrCodeInternalError,
			Message: "could not save nonce",
			Wrapped: err,
		}
	}

	return nonce, nil
}
