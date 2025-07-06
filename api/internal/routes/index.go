package routes

import (
	"net/http"

	"github.com/tasjen/fz/db"
	"github.com/tasjen/fz/internal/handlers"
)

func Routes(db *db.Queries) http.Handler {
	mux := http.NewServeMux()
	handlers := handlers.NewHandlers(db)
	mux.HandleFunc("POST /users", handlers.CreateUserHandler)
	return mux
}
