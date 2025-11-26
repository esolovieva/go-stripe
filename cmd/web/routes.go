package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()                            //multiplexer - request router. Could be also called router
	mux.Get("/virtual-terminal", app.VirtualTerminal) //execute app.VirtualTerminal function after receiving a '[baseUrl]/virtual-terminal' get request

	return mux
}
