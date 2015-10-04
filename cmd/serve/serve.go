package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	app = flag.String("app", "", "Root directory from which to serve the frontend app")
	src = flag.String("src", "", "Root directory from which to serve raw source code")
)

func main() {
	flag.Parse()

	r := mux.NewRouter()

	r.PathPrefix("/src").Methods("GET").
		Handler(http.StripPrefix("/src", http.FileServer(http.Dir(*src))))

	r.PathPrefix("/").Methods("GET").
		Handler(http.FileServer(http.Dir(*app)))

	http.ListenAndServe(":8080", r)
}
