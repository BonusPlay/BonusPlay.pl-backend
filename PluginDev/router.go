package main

import (
	"context"
	"net/http"

	"github.com/BonusPlay/VueHoster/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

type Router struct{}

var server http.Server

func (p Router) Run() (err error) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.ServeFile("dev_files/index.html", w)
	})

	router.Get("/cv_en", func(w http.ResponseWriter, r *http.Request) {
		util.ServeFile("dev_files/cv_en.pdf", w)
	})

	router.Get("/cv_pl", func(w http.ResponseWriter, r *http.Request) {
		util.ServeFile("dev_files/cv_pl.pdf", w)
	})

	router.Get("/github", http.RedirectHandler("https://github.com/BonusPlay", 301).ServeHTTP)
	router.Get("/facebook", http.RedirectHandler("https://facebook.com/BonusPlay3", 301).ServeHTTP)
	router.Get("/discord", http.RedirectHandler("https://discordapp.com/invite/tYk4PW5", 301).ServeHTTP)
	router.Get("/youtube", http.RedirectHandler("https://www.youtube.com/user/adamklis1975", 301).ServeHTTP)
	router.Get("/asktoask", http.RedirectHandler("https://www.youtube.com/watch?v=53zkBvL4ZB4", 301).ServeHTTP)

	// static files
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		util.ServeFile("dev_files/index.html", w)
	})

	log.Info("Starting main on port 3011")
	server = http.Server{
		Addr:    ":3011",
		Handler: router,
	}
	return server.ListenAndServe()
}

func (p Router) Cancel() {
	log.Debug("Main shutting down")
	_ = server.Shutdown(context.Background())
}

// exported plugin
var Plugin Router
