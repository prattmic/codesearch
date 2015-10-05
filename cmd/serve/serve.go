package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/prattmic/codesearch/pkg/search"
)

var (
	app = flag.String("app", "", "Root directory from which to serve the frontend app")
	src = flag.String("src", "", "Root directory from which to serve raw source code")

	index  = flag.String("index", "", "Index to search")
	prefix = flag.String("prefix", "", "Prefix on index paths")
)

var searcher *search.Searcher

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if searcher == nil {
		http.Error(w, "Search unavailable", 500)
		return
	}

	query, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read query: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	results, err := searcher.Search(string(query))
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, fmt.Sprintf("Search error: %v", err), 400)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Unable to write results: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}
}

func main() {
	flag.Parse()

	if *index != "" {
		searcher = search.NewSearcher(*index, *prefix)
	}

	r := mux.NewRouter()

	r.PathPrefix("/src").Methods("GET").
		Handler(http.StripPrefix("/src", http.FileServer(http.Dir(*src))))

	r.PathPrefix("/search").Methods("POST").HandlerFunc(searchHandler)

	r.PathPrefix("/").Methods("GET").
		Handler(http.FileServer(http.Dir(*app)))

	http.ListenAndServe(":8080", r)
}
