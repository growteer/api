package environment

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

func Load() *Environment {
	if err := godotenv.Load(); err != nil {
		slog.Debug("no .env file loaded")
	}
	
	env := new(Environment)

	if err := envconfig.Init(env); err != nil {
		panic(err)
	}

	return env
}

type Environment struct {
	Mongo MongoEnv
	Server ServerEnv
	Token TokenEnv
}

type MongoEnv struct {
	Host string `envconfig:"MONGO_HOST"`
	Port int `envconfig:"MONGO_PORT"`
	User string `envconfig:"MONGO_USER"`
	Password string `envconfig:"MONGO_PASSWORD"`
	DBName string `envconfig:"MONGO_DB_NAME"`
	SSL bool `envconfig:"MONGO_SSL,default=true"`
}

type ServerEnv struct {
	HTTPPort int `envconfig:"HTTP_PORT,default=8080"`
}

type TokenEnv struct {
	JWTSecret string `envconfig:"JWT_SECRET"`
	SessionTTLMinutes int `envconfig:"SESSION_TTL_MINUTES,default=15"`
}
