package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// routes configures the application's HTTP router and returns it.
// It initializes a new Chi router, applies CORS middleware with the
// specified options, and prepares the handler for use by the HTTP server.
//
// The returned handler includes:
//   - CORS configuration allowing HTTP and HTTPS origins
//   - Allowed methods: GET, POST, PUT, DELETE, OPTIONS
//   - Allowed headers such as Accept, Authorization, Content-Type
//   - Disabled credential sharing
//   - A preflight cache duration of 300 seconds (5 minutes)
//
// Additional routes and middleware can be added to the router before
// the server starts handling requests.
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/api/payment-intent", app.GetPaymentIntent)

	return mux
}
