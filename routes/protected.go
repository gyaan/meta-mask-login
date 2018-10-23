package routes

import (
	"github.com/rs/cors"
	"github.com/go-chi/chi"
	"github.com/gyaan/meta-mask-login/middlewares"
	"github.com/gyaan/meta-mask-login/controller"
	 r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func Protected(s *r.Session, cors *cors.Cors) func(r chi.Router)  {

	return func(r chi.Router) {
		r.Use(middlewares.TokenAuthentication)
		r.Use(cors.Handler)
		r.Get("/user/files",controller.GetUserFiles(s))
	}
}