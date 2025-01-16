package graph

import (
	"github.com/growteer/api/infrastructure/environment"
	"github.com/growteer/api/infrastructure/tokens"
	"github.com/growteer/api/internal/authn"
	"github.com/growteer/api/internal/profiles"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	authnService *authn.Service
}

func NewResolver(db *mongo.Database, env *environment.Environment) *Resolver {
	tokenProvider := tokens.NewProvider(env.Token.JWTSecret, env.Token.SessionTTLMinutes, env.Token.RefreshTTLMinutes)

	authnRepo, err := authn.NewRepository(db)
	if err != nil {
		panic(err)
	}

	profileRepo := profiles.NewRepository(db)

	return &Resolver{
		authnService: authn.NewService(authnRepo, tokenProvider, profileRepo),
	}
}