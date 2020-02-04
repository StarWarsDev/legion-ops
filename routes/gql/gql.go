package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/StarWarsDev/legion-ops/internal/gql"
	"github.com/StarWarsDev/legion-ops/internal/gql/resolvers"
	"github.com/gorilla/sessions"
)

type GraphQLHandlers struct {
	store *sessions.CookieStore
}

func New(store *sessions.CookieStore) GraphQLHandlers {
	return GraphQLHandlers{store: store}
}

func (gh *GraphQLHandlers) GraphQLHandler() http.Handler {
	c := gql.Config{
		Resolvers: &resolvers.Resolver{},
	}

	return handler.GraphQL(gql.NewExecutableSchema(c))
}

func (gh *GraphQLHandlers) GraphicalHandler(path string) http.Handler {
	return handler.Playground("Legion Ops GraphQL", path)
}
