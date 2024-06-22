package main

import (
	"errors"
	"fcc-project/cmd/config"
	"fcc-project/internal/models"
	"fcc-project/internal/validator"
	"fmt"
	"net/http"
	"strconv"
)

// Define a createSnippetFormData struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.

// Update our snippetCreateForm struct to include struct tags which tell the
// decoder how to map HTML form values into the different struct fields. So, for
// example, here we're telling the decoder to store the value from the HTML form
// input with the name "title" in the Title field. The struct tag `form:"-"`
// tells the decoder to completely ignore a field during decoding.
type createSnippetFormData struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func home(app *config.Application) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(responseWriter, err)
			return
		}

		// Create an instance of a TemplateData struct holding the data.
		data := app.NewTemplateData(request)
		data.Snippets = snippets
		app.Render(responseWriter, http.StatusOK, "home.html", data)
	}
}

func snippetView(app *config.Application) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		id, err := strconv.Atoi(request.PathValue("id"))
		if err != nil || id < 1 {
			app.NotFound(responseWriter)
			return
		}

		// Use the SnippetModel object's Get method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		snippet, err := app.Snippets.Get(id)

		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(responseWriter)
			} else {
				app.ServerError(responseWriter, err)
			}
			return
		}

		data := app.NewTemplateData(request)
		data.Snippet = snippet

		app.Render(responseWriter, http.StatusOK, "view.html", data)
	}
}

func snippetCreateForm(app *config.Application) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		data := app.NewTemplateData(request)
		// Initialize a new createSnippetForm instance and pass it to the template.
		// Notice how this is also a great opportunity to set any default or
		// 'initial' values for the form --- here we set the initial value forthe snippet expiry to 365 days.
		data.Form = createSnippetFormData{
			Expires: 365,
		}
		app.Render(responseWriter, http.StatusOK, "create.html", data)
	}
}

func snippetCreatePost(app *config.Application) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		var form createSnippetFormData
		err := app.DecodePostForm(request, &form)
		if err != nil {
			app.ClientError(responseWriter, http.StatusBadRequest)
			return
		}

		// Because the Validator type is embedded by the snippetCreateForm struct,
		// we can call CheckField() directly on it to execute our validation checks.
		// CheckField() will add the provided key and error message to the
		// FieldErrors map if the check does not evaluate to true.
		// For example, in the first line here we "check that the form.Title field is not blank".
		// In the second, we "check that the form.Title field has a maximum character length of 100" and so on.
		form.Validator.CheckField(
			validator.NotBlank(form.Title),
			"title",
			"This field cannot be blank",
		)
		form.Validator.CheckField(
			validator.MaxChars(form.Title, 100),
			"title",
			"This field cannnot be more than 100 characters long",
		)
		form.Validator.CheckField(
			validator.NotBlank(form.Content),
			"content",
			"This field cannot be blank",
		)
		form.Validator.CheckField(
			validator.PermittedInt(form.Expires, 1, 7, 365),
			"expires",
			"This field must equal 1, 7 or 365",
		)

		// If there are any validation errors re-display the create.html template,
		// passing in the snippetCreateForm instance as dynamic data in the Form
		// field. Note that we use the HTTP status code 422 Unprocessable Entity
		// when sending the response to indicate that there was a validation error.
		if !form.Valid() {
			data := app.NewTemplateData(request)
			data.Form = form
			app.Render(responseWriter, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.ServerError(responseWriter, err)
			return
		}
		app.SessionManager.Put(request.Context(), "flash", "Snippert successfully created!")
		http.Redirect(responseWriter, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}
