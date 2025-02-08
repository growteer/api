package main

import (
	"github.com/growteer/api/internal/api"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/internal/infrastructure/tokens"
)

func main() {
	env := environment.Load()

	db := mongodb.NewDB(env.Mongo)
	tokenProvider := tokens.NewProvider(env.Token.JWTSecret, env.Token.SessionTTLMinutes, env.Token.RefreshTTLMinutes)

	server := api.NewServer(env.Server, db, tokenProvider)

	server.Start()
}
