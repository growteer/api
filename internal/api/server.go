package api

import (
	"context"
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
	server http.Server
	Router *chi.Mux
	port   int
}

func (s *GQLServer) Start() {
	slog.Info(fmt.Sprintf("connect on port %d for GraphQL playground", s.port))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server unexpectedly shut down", "error", err)
		} else {
			slog.Info("server shut down gracefully")
		}
	}()
}

func (s *GQLServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func NewServer(env environment.Server, db *mongo.Database, tokenProvider authn.TokenProvider) *GQLServer {
	resolver := graphql.NewResolver(db, tokenProvider)

	handler := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))
	handler.SetErrorPresenter(gqlutil.PresentError)
	handler.SetRecoverFunc(gqlutil.Recover)

	handler.AddTransport(transport.Options{})
	handler.AddTransport(transport.GET{})
	handler.AddTransport(transport.POST{})

	handler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	handler.Use(extension.Introspection{})
	handler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router := newRouter(env, handler, tokenProvider)

	return &GQLServer{
		Router: router,
		server: http.Server{Addr: fmt.Sprintf(":%d", env.HTTPPort), Handler: router},
		port:   env.HTTPPort,
	}
}

func newRouter(env environment.Server, handler *handler.Server, tokenProvider authn.TokenProvider) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowedOrigins:   env.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}).Handler)

	router.Use(authn.UserSessionMiddleware(tokenProvider))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", handler)

	return router
}
