package main

import (
	"dream_dictionary/internals"
	_ "embed"
	"log"
	"net/http"
)

type App struct {
	*internals.Data
}

func main() {
	data := &internals.Data{}
	data.Populate()

	app := &App{Data: data}

	mux := http.NewServeMux()
	mux.Handle("/search", http.HandlerFunc(app.SearchIndex))
	mux.Handle("/", http.HandlerFunc(app.Home))
	mux.Handle("/search/content", http.HandlerFunc(app.SearchContent))
	mux.Handle("/search/index", http.HandlerFunc(app.BlogInternalSearch))
	log.Println("Listening on :8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}
