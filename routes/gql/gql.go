package gql

import (
	"net/http"

	"github.com/StarWarsDev/legion-ops/internal/orm"

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

func (gh *GraphQLHandlers) GraphQLHandler(dbORM *orm.ORM) http.Handler {
	c := gql.Config{
		Resolvers: &resolvers.Resolver{
			ORM: dbORM,
		},
	}

	return handler.GraphQL(gql.NewExecutableSchema(c))
}

func (gh *GraphQLHandlers) GraphicalHandler(path string) http.Handler {
	return handler.Playground("Legion Ops GraphQL", path)
}
