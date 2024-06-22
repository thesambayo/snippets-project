package config

import (
	"fcc-project/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

// TemplateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type TemplateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

func HumanDate(date time.Time) string {
	return date.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": HumanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all file paths that
	// match the pattern "./ui/html/pages/*.html". This will essentially gives
	// us a slice of all the file paths for our application 'page' templates
	// like: [ui/html/pages/home.html ui/html/pages/view.html]

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// loop through the page filePaths one-by-one
	for _, page := range pages {
		// extract the file name (like home.html) from the full filePath and assign it to a variable
		name := filepath.Base(page)

		// registering templates functions
		// the template.FuncMap must be registered with the template set before you call the ParseFiles() method
		// this means we have to use template.New() to create an empty template set
		// use the Funcs() method to register the template FuncMap, and then parse the file as normal

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		//ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add the template set to the map, using the name of the page
		// (like 'home.html') as the map key.
		cache[name] = ts
	}

	return cache, nil
}
