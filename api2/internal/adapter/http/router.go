package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tasjen/fz/api-hexa/internal/adapter/config"
	"github.com/tasjen/fz/api-hexa/internal/adapter/http/handler"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.HTTP, userHandler *handler.UserHandler) (*Router, error) {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors.Default())

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
		}
	}
	return &Router{
		Engine: router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
