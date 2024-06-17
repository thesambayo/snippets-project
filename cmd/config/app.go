package config

import (
	"fcc-project/internal/models"
	"html/template"
	"log"

	"github.com/go-playground/form/v4"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
	FormDecoder   *form.Decoder
}
