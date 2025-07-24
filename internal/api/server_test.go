package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/growteer/api/internal/api"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/internal/infrastructure/tokens"
	"github.com/test-go/testify/assert"
)

var config = environment.Load()

func TestNewServer(t *testing.T) {
	var server *api.GQLServer

	db := mongodb.NewDB(config.Mongo)
	tokenProvider := tokens.NewProvider(config.Token.JWTSecret, config.Token.SessionTTLMinutes, config.Token.RefreshTTLMinutes)

	t.Run("create server", func(t *testing.T) {
		server = api.NewServer(config.Server, db, tokenProvider)

		assert.NotNil(t, server)
	})

	t.Run("start server", func(t *testing.T) {
		assert.NotPanics(t, func() { server.Start() })
	})

	t.Run("check that server is online", func(t *testing.T) {
		response, err := http.Get(fmt.Sprintf("http://localhost:%d", config.Server.HTTPPort))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("shutdown", func(t *testing.T) {
		assert.NoError(t, server.Shutdown(context.Background()))
	})
}
