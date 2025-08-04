package main

import (
	"log"
	"net/http"

	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"
	
	"github.com/lincentpega/pcrm/internal/config"
	"github.com/lincentpega/pcrm/internal/handlers"
	"github.com/lincentpega/pcrm/internal/handlers/api"
	"github.com/lincentpega/pcrm/internal/middleware"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/templates"
	_ "github.com/lincentpega/pcrm/docs"
)

// @title Personal CRM API
// @version 1.0
// @description A personal CRM application with REST API endpoints alongside HTMX SSR interface
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	renderer := templates.NewRenderer()
	personRepo := repository.NewPersonRepository(db)
	contactRepo := repository.NewContactRepository(db)
	personHandler := handlers.NewPersonHandler(personRepo, contactRepo, renderer)
	contactHandler := handlers.NewContactHandler(contactRepo, renderer)
	
	personAPI := api.NewPersonAPI(personRepo, contactRepo)
	contactAPI := api.NewContactAPI(contactRepo)

	middlewareChain := alice.New(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
	)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/people", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /people", personHandler.ListPeople)
	mux.HandleFunc("GET /people/new", personHandler.NewPersonForm)
	mux.HandleFunc("POST /people", personHandler.CreatePerson)
	mux.HandleFunc("GET /people/{id}", personHandler.ShowPerson)
	mux.HandleFunc("GET /people/{id}/edit", personHandler.EditPersonForm)
	mux.HandleFunc("GET /people/{id}/edit-inline", personHandler.EditPersonInlineForm)
	mux.HandleFunc("PUT /people/{id}", personHandler.UpdatePerson)
	mux.HandleFunc("DELETE /people/{id}", personHandler.DeletePerson)

	mux.HandleFunc("GET /people/{personId}/contacts/new", contactHandler.NewContactForm)
	mux.HandleFunc("POST /people/{personId}/contacts", contactHandler.CreateContact)
	mux.HandleFunc("GET /people/{personId}/contacts/{contactId}/edit", contactHandler.EditContactForm)
	mux.HandleFunc("PUT /people/{personId}/contacts/{contactId}", contactHandler.UpdateContact)
	mux.HandleFunc("DELETE /people/{personId}/contacts/{contactId}", contactHandler.DeleteContact)

	// REST API routes
	mux.HandleFunc("GET /api/people", personAPI.ListPeople)
	mux.HandleFunc("POST /api/people", personAPI.CreatePerson)
	mux.HandleFunc("GET /api/people/{id}", personAPI.GetPerson)
	mux.HandleFunc("GET /api/people/{id}/full", personAPI.GetPersonWithContacts)
	mux.HandleFunc("PUT /api/people/{id}", personAPI.UpdatePerson)
	mux.HandleFunc("DELETE /api/people/{id}", personAPI.DeletePerson)

	mux.HandleFunc("GET /api/people/{personId}/contacts", contactAPI.ListContactsByPerson)
	mux.HandleFunc("POST /api/people/{personId}/contacts", contactAPI.CreateContact)
	mux.HandleFunc("GET /api/contacts/{id}", contactAPI.GetContact)
	mux.HandleFunc("PUT /api/contacts/{id}", contactAPI.UpdateContact)
	mux.HandleFunc("DELETE /api/contacts/{id}", contactAPI.DeleteContact)
	mux.HandleFunc("GET /api/contact-types", contactAPI.ListContactTypes)

	// Swagger documentation
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	handler := middlewareChain.Then(mux)

	log.Printf("Starting server on %s", cfg.Address())
	if err := http.ListenAndServe(cfg.Address(), handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}