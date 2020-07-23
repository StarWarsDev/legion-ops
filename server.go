package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm"

	"github.com/StarWarsDev/legion-ops/routes/gql"

	"github.com/gorilla/sessions"

	"github.com/StarWarsDev/legion-ops/routes/middlewares"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
		return
	})
}

func StartServer(port, localFilePath string, wait time.Duration, dbORM *orm.ORM) {
	storeSalt := os.Getenv("STORE_SALT")
	if storeSalt == "" {
		storeSalt = "LOCAL_DEV"
	}
	store := sessions.NewCookieStore([]byte(storeSalt))
	gob.Register(map[string]interface{}{})

	middlewareFuncs := middlewares.New(store, dbORM)
	graphqlHandlers := gql.New(store)

	n := negroni.Classic()
	r := mux.NewRouter()

	r.Use(CORS)
	r.Use(middlewareFuncs.Authorize)

	r.Handle("/graphical", graphqlHandlers.GraphicalHandler("/graphql"))
	r.Handle("/graphql", graphqlHandlers.GraphQLHandler(dbORM))

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
