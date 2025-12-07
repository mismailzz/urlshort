package main

import (
	"fmt"
	"net/http"

	"github.com/mismailzz/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8081")
	http.ListenAndServe(":8081", yamlHandler)

	//err := http.ListenAndServe(":8081", mapHandler)
	//fmt.Println("Server ended with error:", err)
}

func defaultMux() *http.ServeMux {
	// ServerMux is an HTTP request multiplexer.
	// Go’s HTTP request router. It decides which
	// handler should run based on the incoming request path.
	// In this case "/", it routes to the hello handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
