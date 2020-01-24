package middlewares

import "github.com/gorilla/sessions"

type MiddlewareFuncs struct {
	store *sessions.CookieStore
}

func New(store *sessions.CookieStore) MiddlewareFuncs {
	return MiddlewareFuncs{store: store}
}
