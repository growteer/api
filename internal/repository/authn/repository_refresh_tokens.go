package authn

import (
	"context"
	"time"

	"github.com/growteer/api/pkg/web3util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daoRefreshToken struct {
	DID       string    `bson:"_id"`
	Token     string    `bson:"token"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r *repository) SaveRefreshToken(ctx context.Context, did *web3util.DID, token string) error {
	newRecord := daoRefreshToken{
		DID:       did.String(),
		Token:     token,
		CreatedAt: time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.nonces.ReplaceOne(ctx, bson.M{"_id": did.String()}, newRecord, opts)

	return err
}

func (r *repository) GetRefreshTokenByDID(ctx context.Context, did *web3util.DID) (string, error) {
	var result daoRefreshToken
	err := r.refreshTokens.FindOne(ctx, bson.M{"_id": did.String()}).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}
