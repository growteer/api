package graphql

import (
	authnApp "github.com/growteer/api/internal/app/authn"
	"github.com/growteer/api/internal/app/profiles"
	authnRepo "github.com/growteer/api/internal/repository/authn"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authnService   *authnApp.Service
	profileService *profiles.Service
}

func NewResolver(db *mongo.Database, tokenProvider authnApp.TokenProvider) *Resolver {
	authnRepo, err := authnRepo.NewRepository(db)
	if err != nil {
		panic(err)
	}

	profileRepo := profiles.NewRepository(db)

	return &Resolver{
		authnService:   authnApp.NewService(authnRepo, tokenProvider, profileRepo),
		profileService: profiles.NewService(profileRepo),
	}
}
