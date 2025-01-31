package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/growteer/api/graph"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/growteer/api/internal/infrastructure/mongodb"
	"github.com/growteer/api/internal/infrastructure/session"
	"github.com/growteer/api/internal/infrastructure/tokens"
	"github.com/growteer/api/pkg/gqlutil"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	env := environment.Load()

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   env.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}).Handler)

	db := mongodb.NewDB(env.Mongo)
	tokenProvider := tokens.NewProvider(env.Token.JWTSecret, env.Token.SessionTTLMinutes, env.Token.RefreshTTLMinutes)

	resolver := graph.NewResolver(db, tokenProvider)
	server := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	server.SetErrorPresenter(gqlutil.PresentError)
	server.SetRecoverFunc(gqlutil.Recover)

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.POST{})

	server.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	server.Use(extension.Introspection{})
	server.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.Use(session.UserSessionMiddleware(tokenProvider))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	port := env.Server.HTTPPort
	slog.Info(fmt.Sprintf("connect on port %d for GraphQL playground", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server unexpectedly shut down", "error", err)
	} else {
		slog.Info("server shut down gracefully")
	}
}
