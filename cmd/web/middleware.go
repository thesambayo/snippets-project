package main

import (
	"fcc-project/cmd/config"
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		responseWriter.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
		responseWriter.Header().Set("X-Frame-Options", "deny")
		responseWriter.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(responseWriter, request)
	})
}

func logRequest(next http.Handler, app *config.Application) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.RequestURI())
		app.InfoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())
		next.ServeHTTP(responseWriter, request)
	})
}

func recoverFromPanic(next http.Handler, app *config.Application) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				responseWriter.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.ServerError(responseWriter, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(responseWriter, request)
	})
}
