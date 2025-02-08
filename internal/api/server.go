package api

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
	"github.com/growteer/api/internal/api/graphql"
	"github.com/growteer/api/internal/api/graphql/gqlutil"
	"github.com/growteer/api/internal/app/authn"
	"github.com/growteer/api/internal/infrastructure/environment"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
	"go.mongodb.org/mongo-driver/mongo"
)

type GQLServer struct {
	*handler.Server
	port   int
	router *chi.Mux
}

func (s *GQLServer) Start() {
	slog.Info(fmt.Sprintf("connect on port %d for GraphQL playground", s.port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server unexpectedly shut down", "error", err)
	} else {
		slog.Info("server shut down gracefully")
	}
}

func NewServer(env environment.ServerEnv, db *mongo.Database, tokenProvider authn.TokenProvider) *GQLServer {
	resolver := graphql.NewResolver(db, tokenProvider)

	server := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
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

	return &GQLServer{
		Server: server,
		port:   env.HTTPPort,
		router: newRouter(env, server, tokenProvider),
	}
}

func newRouter(env environment.ServerEnv, server *handler.Server, tokenProvider authn.TokenProvider) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   env.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}).Handler)

	router.Use(authn.UserSessionMiddleware(tokenProvider))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	return router
}
