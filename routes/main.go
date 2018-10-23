package routes

import (
	"net/http"
	"github.com/go-chi/chi"
	cors2 "github.com/rs/cors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func Router(s *r.Session) http.Handler {
	route:= chi.NewRouter()

	//cors setup
	cors := cors2.New(cors2.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"willDefined"},
		AllowCredentials: true,
		MaxAge:           300,
	})


	route.Group(Protected(s, cors))
	route.Group(Public(s, cors))

	return route
}
