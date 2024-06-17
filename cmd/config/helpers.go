package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// The ServerError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *Application) ServerError(responseWriter http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The ClientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *Application) ClientError(responseWriter http.ResponseWriter, status int) {
	http.Error(responseWriter, http.StatusText(status), status)
}

// NotFound :for consistency, we'll also implement a NotFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) Render(responseWriter http.ResponseWriter, status int, page string, data *TemplateData) {
	// retrieve the appropriate template set from the cache map based on the page name
	// (like 'home.html'). If no entry exists in the cache with the provided name,
	// then create a new error and call the ServerError() helper method and return
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(responseWriter, err)
		return
	}

	// if there is an error in the templates rendering, the user still gets a 200 response + error displayed on template
	// to fix this, we need to make the template render a two-stage process
	// first: 'trial render'; write the template into a buffer
	// if writing into buffer fails: we can then respond to the user with correct error message
	// if writing into buffer works: we can write content of buffer into ResponseWriter

	templateBuffer := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.

	err := ts.ExecuteTemplate(templateBuffer, "base", data)
	// Execute the template set (ts) and write the response body. Again, if there
	// is any error we call the the serverError() helper.
	// content of the "base" template will be used/wrriten as response body
	// which in turn will invoke/contain the other html templates (partials and pages)
	if err != nil {
		app.ServerError(responseWriter, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP provided HTTP status code to http.ResponseWriter.
	responseWriter.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	templateBuffer.WriteTo(responseWriter)
}

func (app *Application) NewTemplateData() *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

// Create a new decodePostForm() helper method. The second parameter here destination,
// is the target destination that we want to decode the form data into.
func (app *Application) DecodePostForm(request *http.Request, destination any) error {
	// Call ParseForm() on the request, in the same way that we did in our createSnippetPost handler
	err := request.ParseForm()
	if err != nil {
		return err
	}
	// Call Decode() on our decoder instance, passing the target destination
	// as the first parameter.
	err = app.FormDecoder.Decode(destination, request.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning the error
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err
	}
	return nil
}
