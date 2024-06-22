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

	mux.Handle(
		"GET /{$}",
		app.SessionManager.LoadAndSave(home(app)),
	)
	mux.Handle(
		"GET /snippet/view/{id}",
		app.SessionManager.LoadAndSave(snippetView(app)),
	)
	mux.Handle(
		"GET /snippet/create",
		app.SessionManager.LoadAndSave(snippetCreateForm(app)),
	)
	mux.Handle(
		"POST /snippet/create",
		app.SessionManager.LoadAndSave(snippetCreatePost(app)),
	)

	return recoverFromPanic(logRequest(secureHeaders(mux), app), app)
}
