package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

var root = flag.String("root", "", "Root directory from which to serve files")

// SourceHandler implements http.Handler, serving the file specified by the URL.
// If served with a prefix, use http.StripPrefix on this Handler.
type SourceHandler struct{}

func (s *SourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	log.Printf("File requested: %s", p)

	location := path.Join(*root, p)
	log.Printf("Serving: %s", location)

	d, err := ioutil.ReadFile(location)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "404 Not Found", 404)
		return
	}

	if err := sourceTemplate.Execute(w, d); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}
}

func main() {
	flag.Parse()

	r := mux.NewRouter()

	r.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	r.PathPrefix("/src").Methods("GET").Handler(http.StripPrefix("/src", &SourceHandler{}))

	http.ListenAndServe(":8080", r)
}