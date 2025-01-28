package web3util

import "github.com/gagliardetto/solana-go"

func VerifySolanaPublicKey(publicKeyBase58 string) error {
	_, err := solana.PublicKeyFromBase58(publicKeyBase58)

	return err
}