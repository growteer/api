package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/growteer/api/internal/api"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/internal/infrastructure/tokens"
)

func main() {
	env := environment.Load()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	db := mongodb.NewDB(env.Mongo)
	tokenProvider := tokens.NewProvider(env.Token.JWTSecret, env.Token.SessionTTLMinutes, env.Token.RefreshTTLMinutes)

	server := api.NewServer(env.Server, db, tokenProvider)

	server.Start()

	<-sig

	slog.Info("shutting down")
}
