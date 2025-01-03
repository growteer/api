package authn

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daoRefreshToken struct {
	Address string `bson:"address"`
	Token string `bson:"token"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r *repository) SaveRefreshToken(ctx context.Context, address, token string) error {
	newRecord := daoRefreshToken{
		Address: address,
		Token: token,
		CreatedAt: time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.nonces.ReplaceOne(ctx, bson.M{"address": address}, newRecord, opts)

	return err
}

func (r *repository) GetRefreshTokenByAddress(ctx context.Context, address string) (string, error) {
	var result daoRefreshToken

	err := r.refreshTokens.FindOne(ctx, bson.M{"address": address}).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}