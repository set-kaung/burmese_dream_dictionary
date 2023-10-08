package main

import (
	"dream_dictionary/internals"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

//go:embed data/DreamDictionary.json
var file_data []byte

type App struct {
	*internals.Data
}

func (app *App) SearchContent(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, "Malformed Request", http.StatusBadRequest)
		return
	}
	query := r.PostFormValue("query")
	response := map[string][]string{}
	response["data"] = SearchContent(app.Data, query)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(response)
	if err != nil {
		http.Error(rw, "server error in searching...", http.StatusInternalServerError)
	}
}

func (app *App) SearchIndex(rw http.ResponseWriter, r *http.Request) {
	// var builder strings.Builder
	index := r.URL.Query().Get("id")
	i, err := strconv.Atoi(index)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	// data := app.Data.DetailMap[i]
	response := map[string][]*internals.IndexSearchCache{}
	response["SearchDetails"] = app.Data.DetailMap[i]
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(rw, "Server Error", 500)
	}
}

func (app *App) home(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	response := map[string][]*internals.BlogHeader{}
	response["Data"] = app.Data.Blogs
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	err := encoder.Encode(response)
	if err != nil {
		http.Error(rw, "Server Error", 500)
	}
}

func main() {
	data := &internals.Data{}
	data.Populate(file_data)
	app := &App{Data: data}

	mux := http.NewServeMux()
	mux.Handle("/search", http.HandlerFunc(app.SearchIndex))
	mux.Handle("/", http.HandlerFunc(app.home))
	mux.Handle("/search/content", http.HandlerFunc(app.SearchContent))
	log.Println("Listening on :6969")
	err := http.ListenAndServe(":6969", mux)
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}
