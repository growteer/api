package graphql

import (
	"github.com/growteer/api/internal/authn"
	"github.com/growteer/api/internal/profiles"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authnService   *authn.Service
	profileService *profiles.Service
}

func NewResolver(db *mongo.Database, tokenProvider authn.TokenProvider) *Resolver {
	authnRepo, err := authn.NewRepository(db)
	if err != nil {
		panic(err)
	}

	profileRepo := profiles.NewRepository(db)

	return &Resolver{
		authnService:   authn.NewService(authnRepo, tokenProvider, profileRepo),
		profileService: profiles.NewService(profileRepo),
	}
}
