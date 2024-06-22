package config

import (
	"fcc-project/internal/models"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Snippets       *models.SnippetModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
