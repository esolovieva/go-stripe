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

// config holds the application configuration settings, including
// server parameters, environment flags, database credentials,
// and Stripe API keys.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

// application bundles together the dependencies required by the server,
// including configuration settings, loggers, template cache, and version
// information.
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

// serve initializes and starts the HTTP server using the application's
// configuration and routing setup. It returns an error if the server
// fails to start or stops unexpectedly.
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting Back-end server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

// main is the entry point of the application. It parses command-line
// flags, loads environment variables, configures loggers, initializes
// the application struct, and starts the HTTP server.
func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application enviornment {development|production|maintenance}")
	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
