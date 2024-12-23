package ethereum

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

const eip191Prefix = "\x19Ethereum Signed Message:\n"

func GetAddressFromSignature(message string, signature []byte) (string, error) {
	messageHash := getMessageHash(message)
	signatureBytes, err := sanitizeSignature(signature)
	if err != nil {
		return "", err
	}

	publicKey, err := crypto.SigToPub(messageHash, signatureBytes)
	if err != nil {
			return "", fmt.Errorf("could not get public key from signature: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*publicKey)

	return recoveredAddr.Hex(), nil
}

func sanitizeSignature(signature []byte) ([]byte, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// Adjust the v value (recovery ID, last byte) from the 27/28 value used in ethereum, to the 0/1 value expected by the crypto package
	signature[64] -= 27

	return signature, nil

}

func getMessageHash(message string) []byte {
	prefixedMessage := eip191Prefix + strconv.Itoa(len(message)) + message
	messageHash := crypto.Keccak256Hash([]byte(prefixedMessage))

	return messageHash.Bytes()
}
