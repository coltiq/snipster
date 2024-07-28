package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/coltiq/snipster/internal/model"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger   *slog.Logger
	snippets *model.SnippetModel
} // Handler Dependencies

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address") // Command-line flag for server port
	dsn := flag.String("dsn", "web:colt@/snipster?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &model.SnippetModel{DB: db},
	} // Config struct containing dependencies

	srv := &http.Server{
		Addr:              *addr,
		Handler:           app.routes(),
		ReadHeaderTimeout: 10 * time.Second,
	} // Server Config

	logger.Info("starting server", slog.String("addr", *addr))
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
