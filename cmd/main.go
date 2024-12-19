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
	"github.com/growteer/api/graph"
	"github.com/growteer/api/infrastructure/environment"
	"github.com/growteer/api/infrastructure/mongodb"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	env := environment.Load()

	db := mongodb.NewDB(env.Mongo)
	resolver := graph.NewResolver(db)
	server := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.POST{})

	server.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	server.Use(extension.Introspection{})
	server.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	port := env.Server.HTTPPort
	slog.Info(fmt.Sprintf("connect on port %d for GraphQL playground", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server unexpectedly shut down", "error", err)
	} else {
		slog.Info("server shut down gracefully")
	}
}
