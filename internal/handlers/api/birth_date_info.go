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

type BirthDateInfoAPI struct {
	repo       *repository.BirthDateInfoRepository
	personRepo *repository.PersonRepository
}

func NewBirthDateInfoAPI(repo *repository.BirthDateInfoRepository, personRepo *repository.PersonRepository) *BirthDateInfoAPI {
	return &BirthDateInfoAPI{
		repo:       repo,
		personRepo: personRepo,
	}
}

// GetBirthDateInfo godoc
// @Summary Get birth date info for a person
// @Description Get the birth date information for a specific person
// @Tags birth-date-info
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 200 {object} BirthDateInfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/people/{personId}/birth-date-info [get]
func (api *BirthDateInfoAPI) GetBirthDateInfo(w http.ResponseWriter, r *http.Request) {
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

	birthDateInfo, err := api.repo.GetByPersonID(personID)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteNotFound(w, "Birth date info not found")
			return
		}
		WriteInternalError(w, "Failed to fetch birth date info")
		return
	}

	response := mappers.BirthDateInfoDomainToResponse(birthDateInfo)
	WriteSuccess(w, response)
}

// UpsertBirthDateInfo godoc
// @Summary Create or update birth date info
// @Description Create or update the birth date information for a specific person
// @Tags birth-date-info
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Param birthDateInfo body BirthDateInfoRequest true "Birth date info data"
// @Success 200 {object} BirthDateInfoResponse "Updated"
// @Success 201 {object} BirthDateInfoResponse "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/birth-date-info [put]
func (api *BirthDateInfoAPI) UpsertBirthDateInfo(w http.ResponseWriter, r *http.Request) {
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

	var req dto.BirthDateInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if err := validators.ValidateBirthDateInfoRequest(&req); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	birthDateInfo := mappers.BirthDateInfoRequestToDomain(personID, &req)

	// Check if birth date info already exists to determine status code
	existingBirthDateInfo, err := api.repo.GetByPersonID(personID)
	isUpdate := err == nil && existingBirthDateInfo != nil

	if err := api.repo.Upsert(birthDateInfo); err != nil {
		WriteInternalError(w, "Failed to save birth date info")
		return
	}

	response := mappers.BirthDateInfoDomainToResponse(birthDateInfo)
	
	if isUpdate {
		WriteSuccess(w, response)
	} else {
		WriteCreated(w, response)
	}
}

// DeleteBirthDateInfo godoc
// @Summary Delete birth date info
// @Description Delete the birth date information for a specific person
// @Tags birth-date-info
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/birth-date-info [delete]
func (api *BirthDateInfoAPI) DeleteBirthDateInfo(w http.ResponseWriter, r *http.Request) {
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
		WriteInternalError(w, "Failed to delete birth date info")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}