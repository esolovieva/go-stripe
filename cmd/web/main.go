package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

type config struct {
	port int
	env  string
	api  string //URL external API
	db   struct {
		dsn string //connection string to db
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port), //Server Port
		Handler:           app.routes(),                        //router, returns chi or mux
		IdleTimeout:       30 * time.Second,                    //connection downtime
		ReadTimeout:       10 * time.Second,                    //request reading time limit
		ReadHeaderTimeout: 5 * time.Second,                     //limiting the time spent reading headlines
		WriteTimeout:      30 * time.Second,                    //response recording restriction
	}

	app.infoLog.Printf("Starting HTTP server in %s mode on the port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	//Reading command line arguments. Example: go run . -port=8080 -env=production
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")
	flag.Parse()

	//Read Stripe key and secret from environmental var
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	//Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Create template cache
	tc := make(map[string]*template.Template)

	//Create application object
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	//Run server
	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
