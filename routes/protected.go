package routes

import (
	"github.com/rs/cors"
	"github.com/go-chi/chi"
	"github.com/gyaan/meta-mask-login/middlewares"
	"github.com/gyaan/meta-mask-login/controller"
	 r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"net/http"
)

func Protected(s *r.Session, cors *cors.Cors) func(r chi.Router)  {

	return func(r chi.Router) {
		r.Use(middlewares.TokenAuthentication)
		r.Use(cors.Handler)
		r.Get("/user/files",controller.GetUserFiles(s))
		r.Post("/user/file",controller.UploadFile(s))
		r.Options("/user/file", func(writer http.ResponseWriter, request *http.Request) {
			return
		})
	}
}