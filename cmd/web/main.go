package main

import (
	"fcc-project/cmd/config"
	"fcc-project/internal/models"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

func main() {
	// Define command-line flags for the address and MySQL DSN string.
	addr := flag.String("addr", ":4400", "HTTP network address")
	dsn := flag.String("dsn", "web:Sambayo@10@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := config.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// connection pool is closed before the main() function exits.
	defer db.Close()

	// initialize a template cache
	templateCache, err := config.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// initialize a decoder instance...
	formDecoder := form.NewDecoder()
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		// Initialize a models.SnippetModel instance and add it to the application dependencies.
		Snippets:       &models.SnippetModel{DB: db},
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(app),
	}

	app.InfoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
