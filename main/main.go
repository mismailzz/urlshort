package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/mismailzz/urlshort"
)

func main() {

	// Flags
	file := flag.String("file", "", "a yaml/json file containing path-url mappings")
	flag.Parse()

	mux := defaultMux()
	var mapHandler http.Handler = mux // default to mux
	if *file != "" {
		handler, err := urlshort.BuildHandlerFromFile(*file, mux /* fallback*/)
		if err != nil {
			panic(err)
		}
		mapHandler = handler
	}

	fmt.Println("Starting the server on :8081")
	err := http.ListenAndServe(":8081", mapHandler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
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
