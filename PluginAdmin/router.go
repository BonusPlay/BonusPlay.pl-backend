package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Router struct{}

var server http.Server

func (p Router) Run() (err error) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// static files
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		if data, err := ioutil.ReadFile("admin_files" + r.URL.Path); err == nil {
			_, _ = w.Write(data)
		} else {
			http.ServeFile(w, r, "admin_files/index.html")
		}
	})

	log.Info("Starting admin on port 3011")
	server = http.Server{
		Addr:    ":1339",
		Handler: router,
	}
	return server.ListenAndServe()
}

func (p Router) Cancel() {
	log.Debug("Main shutting down")
	_ = server.Shutdown(context.Background())
}

// exported plugin
//noinspection GoUnusedGlobalVariable
var Plugin Router
