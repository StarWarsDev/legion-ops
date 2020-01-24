package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

type UserHandlers struct {
	store *sessions.CookieStore
}

func New(store *sessions.CookieStore) UserHandlers {
	return UserHandlers{store: store}
}

func (uh *UserHandlers) ApiMeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := uh.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := session.Values["profile"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
