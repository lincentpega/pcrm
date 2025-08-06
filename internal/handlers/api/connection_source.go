package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/mappers"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/validators"
)

type ConnectionSourceAPI struct {
	repo     *repository.ConnectionSourceRepository
	personRepo *repository.PersonRepository
}

func NewConnectionSourceAPI(repo *repository.ConnectionSourceRepository, personRepo *repository.PersonRepository) *ConnectionSourceAPI {
	return &ConnectionSourceAPI{
		repo:     repo,
		personRepo: personRepo,
	}
}

// GetConnectionSource godoc
// @Summary Get connection source for a person
// @Description Get the connection source information for how we met a specific person
// @Tags connection-sources
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 200 {object} ConnectionSourceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/people/{personId}/connection-source [get]
func (api *ConnectionSourceAPI) GetConnectionSource(w http.ResponseWriter, r *http.Request) {
	personID, err := validators.ValidatePersonID(r.PathValue("personId"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	// Check if person exists
	_, err = api.personRepo.GetByID(personID)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	connectionSource, err := api.repo.GetByPersonID(personID)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteNotFound(w, "Connection source not found")
			return
		}
		WriteInternalError(w, "Failed to fetch connection source")
		return
	}

	response := mappers.ConnectionSourceDomainToResponse(connectionSource)
	WriteSuccess(w, response)
}

// UpsertConnectionSource godoc
// @Summary Create or update connection source
// @Description Create or update the connection source information for how we met a specific person
// @Tags connection-sources
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Param connectionSource body ConnectionSourceRequest true "Connection source data"
// @Success 200 {object} ConnectionSourceResponse "Updated"
// @Success 201 {object} ConnectionSourceResponse "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/connection-source [put]
func (api *ConnectionSourceAPI) UpsertConnectionSource(w http.ResponseWriter, r *http.Request) {
	personID, err := validators.ValidatePersonID(r.PathValue("personId"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	// Check if person exists
	_, err = api.personRepo.GetByID(personID)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	var req dto.ConnectionSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if err := validators.ValidateConnectionSourceRequest(&req); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	connectionSource := mappers.ConnectionSourceRequestToDomain(personID, &req)

	// Check if connection source already exists to determine status code
	existingConnectionSource, err := api.repo.GetByPersonID(personID)
	isUpdate := err == nil && existingConnectionSource != nil

	if err := api.repo.Upsert(connectionSource); err != nil {
		WriteInternalError(w, "Failed to save connection source")
		return
	}

	response := mappers.ConnectionSourceDomainToResponse(connectionSource)
	
	if isUpdate {
		WriteSuccess(w, response)
	} else {
		WriteCreated(w, response)
	}
}

// DeleteConnectionSource godoc
// @Summary Delete connection source
// @Description Delete the connection source information for a specific person
// @Tags connection-sources
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/connection-source [delete]
func (api *ConnectionSourceAPI) DeleteConnectionSource(w http.ResponseWriter, r *http.Request) {
	personID, err := validators.ValidatePersonID(r.PathValue("personId"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	// Check if person exists
	_, err = api.personRepo.GetByID(personID)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	if err := api.repo.Delete(personID); err != nil {
		WriteInternalError(w, "Failed to delete connection source")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}