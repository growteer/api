package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	solana_go "github.com/gagliardetto/solana-go"
)


func VerifyPublicKey(publicKeyBase58 string) error {
	_, err := solana_go.PublicKeyFromBase58(publicKeyBase58)

	return err
}

func VerifySignature(message, signatureBase64, publicKeyBase58 string) error {
	publicKey, err := solana_go.PublicKeyFromBase58(publicKeyBase58)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}

	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature length: %d", len(signature))
	}
	
	isValid := ed25519.Verify(publicKey[:], []byte(message), signature)
	if !isValid {
		return fmt.Errorf("invalid signature for message and public key")
	}

	return nil
}
