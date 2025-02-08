package authn

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const db_collection_nonces = "nonces"
const db_collection_refresh_tokens = "refresh_tokens"

type repository struct {
	nonces        *mongo.Collection
	refreshTokens *mongo.Collection
}

func NewRepository(db *mongo.Database) (*repository, error) {
	repo := &repository{
		nonces:        db.Collection(db_collection_nonces),
		refreshTokens: db.Collection(db_collection_refresh_tokens),
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "createdAt", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(10 * 60), // Auto-delete Nonces after 10 minutes
	}

	_, err := repo.nonces.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
