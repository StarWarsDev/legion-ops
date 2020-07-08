package middlewares

import (
	"github.com/StarWarsDev/legion-ops/internal/orm"
	"github.com/gorilla/sessions"
)

type MiddlewareFuncs struct {
	store *sessions.CookieStore
	dbORM *orm.ORM
}

func New(store *sessions.CookieStore, dbORM *orm.ORM) MiddlewareFuncs {
	return MiddlewareFuncs{store: store, dbORM: dbORM}
}
