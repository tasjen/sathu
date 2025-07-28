package routes

import (
	"net/http"

	"github.com/tasjen/fz/api/internal/handlers"
)

type Router struct {
	handlers *handlers.Handlers
}

func NewRouter(handlers *handlers.Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) RegisterRoutes() http.Handler {
	router := http.NewServeMux()
	router.Handle("/users/", http.StripPrefix("/users", r.userRouter()))
	return router
}
