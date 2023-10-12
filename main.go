package main

import (
	"dream_dictionary/internals"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	*internals.Data
}

type requestJSON struct {
	Query string
}

func (app *App) SearchContent(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	encoder := json.NewEncoder(rw)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Println("Unrecognized data:")
		http.Error(rw, "Malformed Request", http.StatusBadRequest)
		return
	}
	fmt.Println(r.PostForm)
	response := map[string][]string{}
	response["data"] = SearchContent(app.Data, r.PostForm.Get("query"))
	rw.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "application/json")
	err = encoder.Encode(response)
	if err != nil {
		http.Error(rw, "server error in searching...", http.StatusInternalServerError)
	}
	log.Printf("Took %s to respond from /search/content.\n", time.Since(start).String())
}

func (app *App) SearchIndex(rw http.ResponseWriter, r *http.Request) {
	// var builder strings.Builder
	start := time.Now()
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
	log.Printf("took %s to respond from /search.\n", time.Since(start).String())
}

func (app *App) Home(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	encoder := json.NewEncoder(rw)
	response := map[string][]*internals.BlogHeader{}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	response["Data"] = app.Data.Blogs
	err := encoder.Encode(response)
	if err != nil {
		http.Error(rw, "Server Error", 500)
	}
	log.Printf("took %s to respond from root.\n", time.Since(start).String())
}

func main() {
	data := &internals.Data{}
	data.Populate()
	app := &App{Data: data}

	mux := http.NewServeMux()
	mux.Handle("/search", http.HandlerFunc(app.SearchIndex))
	mux.Handle("/", http.HandlerFunc(app.Home))
	mux.Handle("/search/content", http.HandlerFunc(app.SearchContent))
	log.Println("Listening on :6969")
	err := http.ListenAndServe(":6969", mux)
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}
