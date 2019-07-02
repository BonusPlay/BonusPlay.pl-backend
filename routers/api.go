package routers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func ApiRouter() (ret *chi.Mux) {
	ret = chi.NewRouter()

	ret.Use(middleware.RequestID)
	ret.Use(middleware.RealIP)
	ret.Use(middleware.Logger)
	ret.Use(middleware.Recoverer)

	ret.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome to bonusplay.pl API"))
	})

	return
}
