package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

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

	results, err := searcher.Search(search.Options{
		Regexp:  string(query),
		Context: 2,
	})
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, fmt.Sprintf("Search error: %v", err), 400)
		return
	}

	// Serve no more than 10 files' results.
	total := len(results)
	if total > 10 {
		results = results[:10]
	}

	log.Printf("Query: %q, Total results: %d, Serving: %+v", query, total, results)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("Unable to write results: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}
}

// SingleFile always serves the same file.
type SingleFile struct {
	Path string
}

// ServeHTTP serves the file at Path.
func (f *SingleFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(f.Path)
	if err != nil {
		log.Printf("Error opening %v: %v", f.Path, err)
		http.Error(w, "File not found", 404)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error copying %v: %v", f.Path, err)
		http.Error(w, "Internal Server Error", 500)
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

	r.PathPrefix("/api/search").Methods("POST").HandlerFunc(searchHandler)

	// Single-page app URLs are always served with the index page.
	r.PathPrefix("/file/").Methods("GET").
		Handler(&SingleFile{path.Join(*app, "index.html")})
	r.Path("/search").Methods("GET").
		Handler(&SingleFile{path.Join(*app, "index.html")})

	r.PathPrefix("/").Methods("GET").
		Handler(http.FileServer(http.Dir(*app)))

	http.ListenAndServe(":8080", r)
}
