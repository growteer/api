package environment

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

func Load() Environment {
	if err := godotenv.Load(); err != nil {
		slog.Debug("no .env file loaded")
	}

	env := Environment{}

	if err := envconfig.Init(&env); err != nil {
		panic(err)
	}

	return env
}

type Environment struct {
	Mongo  Mongo
	Server Server
	Token  Token
}

type Mongo struct {
	URI      string `envconfig:"MONGODB_URI"`
	Database string `envconfig:"MONGODB_DATABASE"`
}

type Server struct {
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS"`
	HTTPPort       int      `envconfig:"HTTP_PORT,default=8080"`
}

type Token struct {
	JWTSecret         string `envconfig:"JWT_SECRET"`
	SessionTTLMinutes int    `envconfig:"SESSION_TTL_MINUTES,default=15"`
	RefreshTTLMinutes int    `envconfig:"REFRESH_TTL_MINUTES,default=10080"` // Default = One Week
}
