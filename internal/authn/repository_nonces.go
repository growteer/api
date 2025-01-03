package authn

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daoNonce struct {
	Address string `bson:"address"`
	Nonce string `bson:"nonce"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r *repository) GetNonceByAddress(ctx context.Context, address string) (string, error) {
	if !IsValidEthereumAddress(address) {
		return "", fmt.Errorf("invalid ethereum address passed to GetNonceByAddress: %s", address)
	}

	var result daoNonce
	err := r.nonces.FindOne(ctx, bson.M{"address": address}).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Nonce, nil
}

func (r *repository) SaveNonce(ctx context.Context, address, nonce string) error {
	if !IsValidEthereumAddress(address) {
		return fmt.Errorf("invalid ethereum address passed to SaveNonce: %s", address)
	}

	newRecord := daoNonce{
		Address: address,
		Nonce: nonce,
		CreatedAt: time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.nonces.ReplaceOne(ctx, bson.M{"address": address}, newRecord, opts)

	return err
}