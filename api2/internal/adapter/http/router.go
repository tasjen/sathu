package http

import "github.com/tasjen/fz/api-hexa/internal/adapter/http/handler"

type Router struct {
	userHandler *handler.UserHandler
}

func NewRouter(userHandler *handler.UserHandler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

func (r *Router) RegisterRoutes() {
}
