package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/growteer/api/internal/api"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/internal/infrastructure/tokens"
)

func configureLogging() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	configureLogging()

	env := environment.Load()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	db := mongodb.NewDB(env.Mongo)
	tokenProvider := tokens.NewProvider(env.Token.JWTSecret, env.Token.SessionTTLMinutes, env.Token.RefreshTTLMinutes)

	server := api.NewServer(env.Server, db, tokenProvider)

	server.Start()

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("unable to shutdown gracefully", "err", err)

		os.Exit(1)
	}

	slog.Info("shutting down...")
}
