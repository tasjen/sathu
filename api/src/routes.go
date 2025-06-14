package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.helloJson)
	mux.HandleFunc("POST /v1/upload", app.upload)
	return mux
}
