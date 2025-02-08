package authn

import (
	"context"
	"time"

	"github.com/growteer/api/pkg/web3util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daoNonce struct {
	DID string `bson:"_id"`
	Nonce string `bson:"nonce"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r *repository) GetNonceByDID(ctx context.Context, did *web3util.DID) (string, error) {
	var result daoNonce
	err := r.nonces.FindOne(ctx, bson.M{"_id": did.String()}).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Nonce, nil
}

func (r *repository) SaveNonce(ctx context.Context, did *web3util.DID, nonce string) error {
	newRecord := daoNonce{
		DID: did.String(),
		Nonce: nonce,
		CreatedAt: time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.nonces.ReplaceOne(ctx, bson.M{"_id": did.String()}, newRecord, opts)

	return err
}