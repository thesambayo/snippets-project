package main

import (
	"fcc-project/cmd/config"
	"net/http"
)

// The routes() method returns a servemux containing our application routes.
func routes(app *config.Application) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// any routes that matches /static/a/b/...
	mux.Handle("GET /static/{filePath...}", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home(app))
	mux.HandleFunc("GET /snippet/view/{id}", snippetView(app))
	mux.HandleFunc("GET /snippet/create", snippetCreateForm(app))
	mux.HandleFunc("POST /snippet/create", snippetCreatePost(app))

	return recoverFromPanic(logRequest(secureHeaders(mux), app), app)
}
