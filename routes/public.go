package routes

import (
	"github.com/rs/cors"
	"github.com/go-chi/chi"
	"github.com/gyaan/meta-mask-login/controller"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"net/http"
)

func Public(s *r.Session, cors *cors.Cors) func(r chi.Router)  {

	return func(r chi.Router) {

		//add middleware using r.Use
		//no middle as of now
		r.Use(cors.Handler)
		r.Get("/user/{publicAddress}", controller.GetUserDetails(s))
        r.Post("/user",controller.CreateUser(s))

		//cros origin issue
		r.Options("/authenticate", func(writer http.ResponseWriter, request *http.Request) {
			return
		})
		r.Post("/authenticate",controller.Auth(s))
	}
}