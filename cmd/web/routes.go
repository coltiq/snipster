package main

import (
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) // Server static files

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))                        // Restrict this route to exact matches on / only
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.snippetsView)) // Display a specific snippet w/ {id} wildcard
	mux.Handle("GET /users/signup", dynamic.ThenFunc(app.userSignup))         // Display a form for signing up a new user
	mux.Handle("POST /users/signup", dynamic.ThenFunc(app.userSignupPost))    // Create a new user
	mux.Handle("GET /users/login", dynamic.ThenFunc(app.userLogin))           // Display a form for logging in a user
	mux.Handle("POST /users/login", dynamic.ThenFunc(app.userLoginPost))      // Authenticate and login the user

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippets/create", protected.ThenFunc(app.snippetsCreate))      // Disply a form to create a new snippet w/ GET restriction
	mux.Handle("POST /snippets/create", protected.ThenFunc(app.snippetsCreatePost)) // Save a new snippet w/ POST restriction
	mux.Handle("POST /users/logout", protected.ThenFunc(app.userLogoutPost))        // Logout the user

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
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
