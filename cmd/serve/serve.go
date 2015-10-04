package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

var root = flag.String("root", "", "Root directory from which to serve files")

// SourceHandler serves the requested source file.
func SourceHandler(w http.ResponseWriter, r *http.Request) {
	// This is served at /src, so strip that off to get the requested file.
	p := strings.TrimPrefix(r.URL.Path, "/src")
	p = path.Clean(p)
	log.Printf("File requested: %s", p)

	location := path.Join(*root, p)
	log.Printf("Serving: %s", location)

	f, err := os.Open(location)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "404 Not Found", 404)
		return
	}

	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Error copying file: %v", err)
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

	r.PathPrefix("/src").Methods("GET").HandlerFunc(SourceHandler)

	http.ListenAndServe(":8080", r)
}
