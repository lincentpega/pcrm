package main

import (
	"log"
	"net/http"

	"github.com/justinas/alice"
	"github.com/lincentpega/pcrm/internal/config"
	"github.com/lincentpega/pcrm/internal/handlers"
	"github.com/lincentpega/pcrm/internal/middleware"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/templates"
)

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

	handler := middlewareChain.Then(mux)

	log.Printf("Starting server on %s", cfg.Address())
	if err := http.ListenAndServe(cfg.Address(), handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}