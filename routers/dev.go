package routers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"os"
	"path/filepath"
)

func DevRouter() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ServeFile("dev_files/index.html", w)
	})

	r.Get("/cv_en", func(w http.ResponseWriter, r *http.Request) {
		ServeFile("dev_files/cv_en.pdf", w)
	})

	r.Get("/cv_pl", func(w http.ResponseWriter, r *http.Request) {
		ServeFile("dev_files/cv_pl.pdf", w)
	})

	r.Get("/github", http.RedirectHandler("https://github.com/BonusPlay", 301).ServeHTTP)
	r.Get("/facebook", http.RedirectHandler("https://facebook.com/BonusPlay3", 301).ServeHTTP)
	r.Get("/discord", http.RedirectHandler("https://discordapp.com/invite/tYk4PW5", 301).ServeHTTP)
	r.Get("/youtube", http.RedirectHandler("https://www.youtube.com/user/adamklis1975", 301).ServeHTTP)

	// static files
	workDir, _ := os.Getwd()
	r.Get("/*", NoDirListingHandler(http.FileServer(http.Dir(filepath.Join(workDir, "dev_files")))).ServeHTTP)

	return
}