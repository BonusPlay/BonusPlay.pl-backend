package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ServeFile(name string, w http.ResponseWriter) {
	if data, err := ioutil.ReadFile(name); err == nil {
		_, _ = w.Write(data)
	} else {
		log.Fatal(err)
	}
}

// TODO: shouldn't redirect if path points to valid folder name
func NoDirListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
/*
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	log.Println("prealles", path)

	if path != "/" && path[len(path)-1] != '/' {
		//r.Get(path, http.NotFoundHandler().ServeHTTP)
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		log.Println("preadd", path)
		path += "/"
		log.Println("postadd", path)
	}

	// path += "*"

//	fmt.Println(path)
//	if strings.HasSuffix(path, "/") {
//		fmt.Println("prelast", path)
//		r.Get(path, http.NotFoundHandler().ServeHTTP)
//		return
//	}

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("last", path)
		fs.ServeHTTP(w, r)
	}))
}
*/
