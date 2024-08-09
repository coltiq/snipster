package main

import (
	"net/http"

	"github.com/coltiq/snipster/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files)) // Server static files

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
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
