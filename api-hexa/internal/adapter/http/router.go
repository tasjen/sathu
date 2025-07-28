package http

import "github.com/tasjen/fz/api-hexa/internal/adapter/http/handler"

type Router struct {
	userHandler *handler.UserHandler
}
