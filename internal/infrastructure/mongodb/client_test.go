package mongodb_test

import (
	"testing"

	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/test-go/testify/assert"
)

var config = environment.Load()

func TestNewDB(t *testing.T) {
	t.Parallel()

	t.Run("successful connection", func(t *testing.T) {
		database := mongodb.NewDB(config.Mongo)

		assert.NotNil(t, database)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		assert.Panics(t, func() {
			mongodb.NewDB(environment.Mongo{URI: "http://localhost:27017"})
		})
	})

	t.Run("no connection", func(t *testing.T) {
		assert.Panics(t, func() {
			mongodb.NewDB(environment.Mongo{URI: "mongodb://100.100.2.3:11111"})
		})
	})
}
