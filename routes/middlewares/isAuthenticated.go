package middlewares

import (
	"errors"
	"net/http"
)

func (f *MiddlewareFuncs) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := f.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		http.Error(w, errors.New("authentication required").Error(), http.StatusUnauthorized)
		return
	} else {
		next(w, r)
	}
}
