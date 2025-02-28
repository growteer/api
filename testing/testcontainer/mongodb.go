package testcontainer

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "growteer_test"
	containerImage = "mongo:latest"
	defaultPort    = "27017/tcp"
	password       = "mongo"
	ssl            = false
	username       = "mongo"
)

func runContainer() (testcontainers.Container, error) {
	return mongodb.Run(
		context.Background(),
		containerImage,
		mongodb.WithUsername(username),
		mongodb.WithPassword(password),
		testcontainers.WithWaitStrategy(wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort(defaultPort),
		)))
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

	mongoEnv = environment.MongoEnv{
		Host:     host,
		Port:     port.Int(),
		User:     username,
		Password: password,
		DBName:   dbName,
		SSL:      ssl,
	}

	require.NoError(t, waitForReadiness(mongoEnv))

	return mongoEnv, shutdown
}

func waitForReadiness(env environment.MongoEnv) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", env.User, env.Password, env.Host, strconv.Itoa(env.Port)))
	var client *mongo.Client

	for i := 0; i < 30; i++ {
		client, err = mongo.Connect(context.Background(), clientOptions)

		if err == nil {
			err = client.Ping(context.Background(), nil)

			if err == nil {
				break
			}
		}

		time.Sleep(1 * time.Second)
	}

	return
}
