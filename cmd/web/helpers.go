package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (cfg *apiConfig) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	cfg.errorLog.Output(2, trace) // Allows for outputting where the error originally occured

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (cfg *apiConfig) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (cfg *apiConfig) notFound(w http.ResponseWriter) {
	cfg.clientError(w, http.StatusNotFound)
}
