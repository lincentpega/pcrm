package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lincentpega/pcrm/docs"
	"github.com/lincentpega/pcrm/internal/config"
	"github.com/lincentpega/pcrm/internal/handlers/api"
	"github.com/lincentpega/pcrm/internal/middleware"
	"github.com/lincentpega/pcrm/internal/repository"
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

	personRepo := repository.NewPersonRepository(db)
	contactRepo := repository.NewContactRepository(db)

	personAPI := api.NewPersonAPI(personRepo, contactRepo)
	contactAPI := api.NewContactAPI(contactRepo)

	middlewareChain := alice.New(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/people", http.StatusSeeOther)
	})

	// REST API routes
	mux.HandleFunc("GET /api/people", personAPI.ListPeople)
	mux.HandleFunc("POST /api/people", personAPI.CreatePerson)
	mux.HandleFunc("GET /api/people/{id}", personAPI.GetPerson)
	mux.HandleFunc("GET /api/people/{id}/full", personAPI.GetPersonFullInfo)
	mux.HandleFunc("PUT /api/people/{id}", personAPI.UpdatePerson)
	mux.HandleFunc("PUT /api/people/{id}/birthdate", personAPI.UpsertPersonBirthdate)
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

	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
