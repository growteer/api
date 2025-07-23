package graphql

import (
	authnApp "github.com/growteer/api/internal/app/authn"
	profilesApp "github.com/growteer/api/internal/app/profiles"
	authnRepo "github.com/growteer/api/internal/repository/authn"
	profilesRepo "github.com/growteer/api/internal/repository/profiles"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authnService   *authnApp.Service
	profileService *profilesApp.Service
}

func NewResolver(db *mongo.Database, tokenProvider authnApp.TokenProvider) *Resolver {
	authnRepo, err := authnRepo.NewRepository(db)
	if err != nil {
		panic(err)
	}

	profileRepo := profilesRepo.NewRepository(db)

	return &Resolver{
		authnService:   authnApp.NewService(authnRepo, tokenProvider, profileRepo),
		profileService: profilesApp.NewService(profileRepo),
	}
}
