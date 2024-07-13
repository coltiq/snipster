package main

import "net/http"

func (cfg *apiConfig) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./assets/static/")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) // Server static files

	mux.HandleFunc("GET /{$}", cfg.home)                            // Restrict this route to exact matches on / only
	mux.HandleFunc("GET /snippets/view/{id}", cfg.snippetsView)     // Display a specific snippet w/ {id} wildcard
	mux.HandleFunc("GET /snippets/create", cfg.snippetsCreate)      // Disply a form to create a new snippet w/ GET restriction
	mux.HandleFunc("POST /snippets/create", cfg.snippetsCreatePost) // Save a new snippet w/ POST restriction

	return mux
}
