package main

import (
	"net/http"

	"github.com/tasjen/fz/internal/handlers"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	handlers := handlers.NewHandlers(app.db)
	mux.HandleFunc("GET /", handlers.CreateUserHandler)
	return mux
}
