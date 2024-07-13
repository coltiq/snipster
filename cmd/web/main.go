package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address") // Command-line flag for server port
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./assets/static/")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) // Server static files

	mux.HandleFunc("GET /{$}", home)                            // Restrict this route to exact matches on / only
	mux.HandleFunc("GET /snippets/view/{id}", snippetsView)     // Display a specific snippet w/ {id} wildcard
	mux.HandleFunc("GET /snippets/create", snippetsCreate)      // Disply a form to create a new snippet w/ GET restriction
	mux.HandleFunc("POST /snippets/create", snippetsCreatePost) // Save a new snippet w/ POST restriction

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
