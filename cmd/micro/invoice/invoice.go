package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	smtp struct {
		host       string
		port       int
		username   string
		password   string
		encryption string
	}
	frontend string
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func (app *application) serve() error {

	//TODO if we are in a Go stack this microservice should really use gRPC
	// or just RPC of not stack is not in Go
	// it will skip the whole JSON thing and be much more efficient

	//TODO in prod we would probably this microservice would have its own database
	// it would also have additional functionality such as invoice retrieval

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting Invoice Microservice on port %d\n", app.config.port))

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 5000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.smtp.host, "smtphost", "localhost", "smtp host")
	flag.IntVar(&cfg.smtp.port, "smtpport", 1025, "smtp port")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "", "smtp user")
	flag.StringVar(&cfg.smtp.password, "smtppass", "", "smtp password")
	flag.StringVar(&cfg.smtp.encryption, "smtpencryption", "none", "smtp encryption")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to front end")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	app.CreateDirIfNotExists("./invoices")

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
