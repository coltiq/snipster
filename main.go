package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snipster"))
}

func snippetsView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific snippet for id: %d...", id)
	w.Write([]byte(msg))
}

func snippetsCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetsCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)                            // Restrict this route to exact matches on / only
	mux.HandleFunc("GET /snippets/view/{id}", snippetsView)     // Display a specific snippet w/ {id} wildcard
	mux.HandleFunc("GET /snippets/create", snippetsCreate)      // Disply a form to create a new snippet w/ GET restriction
	mux.HandleFunc("POST /snippets/create", snippetsCreatePost) // Save a new snippet w/ POST restriction

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
