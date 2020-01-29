package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gqlhandler "github.com/99designs/gqlgen/handler"
	"github.com/StarWarsDev/legion-ops/routes/login"

	"github.com/StarWarsDev/legion-ops/routes/logout"

	"github.com/StarWarsDev/legion-ops/routes/callback"

	"github.com/StarWarsDev/legion-ops/routes/spa"

	"github.com/gorilla/sessions"

	"github.com/StarWarsDev/legion-ops/routes/middlewares"

	"github.com/StarWarsDev/legion-ops/routes/user"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func StartServer(port, localFilePath string, wait time.Duration) {
	storeSalt := os.Getenv("STORE_SALT")
	if storeSalt == "" {
		storeSalt = "LOCAL_DEV"
	}
	store := sessions.NewCookieStore([]byte(storeSalt))
	gob.Register(map[string]interface{}{})

	middlewareFuncs := middlewares.New(store)
	callbackHandlers := callback.New(store)
	loginHandlers := login.New(store)
	userHandlers := user.New(store)

	n := negroni.Classic()
	r := mux.NewRouter()

	r.HandleFunc("/login", loginHandlers.HandleLogin)
	r.HandleFunc("/logout", logout.Handler)
	r.HandleFunc("/callback", callbackHandlers.HandleCallback)

	r.Handle("/graphical", gqlhandler.Playground("GraphQL playground", "/graphql"))
	r.Handle("/graphql", gqlhandler.GraphQL(NewExecutableSchema(Config{
		Resolvers: &Resolver{},
	})))

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.PathPrefix("/api/me").Handler(n.With(
		negroni.HandlerFunc(middlewareFuncs.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(userHandlers.ApiMeHandler)),
	))

	// this MUST be the final handler
	spaHandler := spa.SPAHandler{
		StaticPath: localFilePath,
		IndexPath:  "index.html",
	}
	r.PathPrefix("/").Handler(spaHandler)

	n.UseHandler(r)
	srv := &http.Server{
		Handler:      n,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
