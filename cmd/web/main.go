package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	infoLog  *log.Logger
	errorLog *log.Logger
} // Handler Dependencies

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address") // Command-line flag for server port
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ltime)                  // Log information (ie "Serving Starting...")
	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ltime|log.Lshortfile) // Log Errors

	apiCfg := &apiConfig{
		infoLog:  infoLog,
		errorLog: errorLog,
	} // Config struct containing dependencies

	srv := &http.Server{
		Addr:              *addr,
		ErrorLog:          errorLog,
		Handler:           apiCfg.routes(),
		ReadHeaderTimeout: 10 * time.Second,
	} // Server Config

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
