package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/growteer/api/internal/infrastructure/environment"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDB(env environment.Mongo) *mongo.Database {
	clientOptions := options.Client().ApplyURI(env.URI).SetTimeout(5 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(fmt.Errorf("failed to connect to mongodb: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(fmt.Errorf("failed to ping mongodb: %w", err))
	}

	slog.Info("connected to mongodb")

	return client.Database(env.Database)
}
