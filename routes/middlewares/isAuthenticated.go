package middlewares

import (
	"encoding/json"
	"net/http"
)

func (f *MiddlewareFuncs) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := f.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		body, _ := json.Marshal(map[string]string{"error": "authentication required"})
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(body), http.StatusUnauthorized)
		return
	} else {
		next(w, r)
	}
}
