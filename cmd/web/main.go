package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
