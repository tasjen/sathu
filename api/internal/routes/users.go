package routes

import "net/http"

func (r *Router) userRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", r.handlers.CreateUserHandler)
	return mux
}
