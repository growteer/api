package nonce

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const db_collection = "nonces"

type dao struct {
	Address string `bson:"address"`
	Nonce string `bson:"nonce"`
	CreatedAt time.Time `bson:"createdAt"`
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) (*repository, error) {
	repo := &repository{ collection: db.Collection(db_collection)}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "createdAt", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(10*60), // Auto-delete Nonces after 10 minutes
	}

	_, err := repo.collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *repository) Save(ctx context.Context, address, nonce string) error {
	newRecord := dao{
		Address: address,
		Nonce: nonce,
		CreatedAt: time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"address": address}, newRecord, opts)

	return err
}
