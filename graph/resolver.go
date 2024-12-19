package graph

import (
	"github.com/growteer/api/internal/nonce"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	nonceService *nonce.Service
}

func NewResolver(db *mongo.Database) *Resolver {
	nonceRepo, err := nonce.NewRepository(db)
	if err != nil {
		panic(err)
	}

	return &Resolver{
		nonceService: nonce.NewService(nonceRepo),
	}
}