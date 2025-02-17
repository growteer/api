package testcontainer

import (
	"context"
	"testing"

	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

const (
	dbName         = "growteer_test"
	containerImage = "mongo:6"
	defaultPort    = "27017/tcp"
	password       = "mongo"
	ssl            = false
	username       = "mongo"
)

func runContainer() (testcontainers.Container, error) {
	return mongodb.Run(context.Background(), containerImage, mongodb.WithUsername(username), mongodb.WithPassword(password))
}

func StartMongoAndGetDetails(t *testing.T) (mongoEnv environment.MongoEnv, shutdown func()) {
	container, err := runContainer()
	require.NoError(t, err)

	shutdown = func() {
		err := container.Terminate(context.Background())
		if err != nil {
			t.Fatal()
		}
	}

	host, err := container.Host(context.Background())
	require.NoError(t, err)

	port, err := container.MappedPort(context.Background(), defaultPort)
	if err != nil {
		shutdown()
		require.Fail(t, err.Error())
	}

	return environment.MongoEnv{
		Host:     host,
		Port:     port.Int(),
		User:     username,
		Password: password,
		DBName:   dbName,
		SSL:      ssl,
	}, shutdown
}
