package api

import (
	"encoding/json"
	"net/http"

	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/mappers"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/validators"
)

type PersonAPI struct {
	repo        *repository.PersonRepository
	contactRepo *repository.ContactRepository
}

func NewPersonAPI(repo *repository.PersonRepository, contactRepo *repository.ContactRepository) *PersonAPI {
	return &PersonAPI{
		repo:        repo,
		contactRepo: contactRepo,
	}
}

// ListPeople godoc
// @Summary List people with pagination
// @Description Get a paginated list of all people in the CRM
// @Tags people
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} PaginatedResponse[PersonInfoResponse]
// @Failure 500 {object} ErrorResponse
// @Router /api/people [get]
func (api *PersonAPI) ListPeople(w http.ResponseWriter, r *http.Request) {
	page, limit := validators.ParsePaginationParams(
		r.URL.Query().Get("page"),
		r.URL.Query().Get("limit"),
	)

	people, err := api.repo.GetPaginated(page, limit)
	if err != nil {
		WriteInternalError(w, "Failed to fetch people")
		return
	}

	totalCount, err := api.repo.GetTotalCount()
	if err != nil {
		WriteInternalError(w, "Failed to get total count")
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	response := make([]dto.PersonInfoResponse, len(people))
	for i, person := range people {
		response[i] = mappers.PersonDomainToResponse(&person)
	}

	WritePaginated(w, response, page, totalPages, totalCount)
}

// GetPerson godoc
// @Summary Get a person by ID
// @Description Get detailed information about a specific person
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} PersonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/people/{id} [get]
func (api *PersonAPI) GetPerson(w http.ResponseWriter, r *http.Request) {
	id, err := validators.ValidatePersonID(r.PathValue("id"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	person, err := api.repo.GetByID(id)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	response := mappers.PersonDomainToResponse(person)
	WriteSuccess(w, response)
}

// GetPersonFullInfo godoc
// @Summary Get complete person information
// @Description Get detailed information about a person including all related data (contacts, etc.)
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} PersonWithContactsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{id}/full [get]
func (api *PersonAPI) GetPersonFullInfo(w http.ResponseWriter, r *http.Request) {
	id, err := validators.ValidatePersonID(r.PathValue("id"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	person, err := api.repo.GetByID(id)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	contacts, err := api.contactRepo.GetByPersonID(id)
	if err != nil {
		WriteInternalError(w, "Failed to fetch contacts")
		return
	}

	response := mappers.PersonWithContactsDomainToResponse(person, contacts)
	WriteSuccess(w, response)
}

// CreatePerson godoc
// @Summary Create a new person
// @Description Create a new person in the CRM system
// @Tags people
// @Accept json
// @Produce json
// @Param person body PersonUpsertRequest true "Person data"
// @Success 201 {object} PersonInfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people [post]
func (api *PersonAPI) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req dto.PersonUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if err := validators.ValidatePersonUpsertRequest(&req); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	person := mappers.PersonUpsertRequestToDomain(&req)
	if err := api.repo.Create(person); err != nil {
		WriteInternalError(w, "Failed to create person")
		return
	}

	response := mappers.PersonDomainToResponse(person)
	WriteCreated(w, response)
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Update an existing person's information
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body PersonUpsertRequest true "Updated person data"
// @Success 200 {object} PersonInfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{id} [put]
func (api *PersonAPI) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id, err := validators.ValidatePersonID(r.PathValue("id"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	var req dto.PersonUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if err := validators.ValidatePersonUpsertRequest(&req); err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	person := mappers.PersonUpsertRequestToDomain(&req)
	person.ID = id

	if err := api.repo.Update(person); err != nil {
		WriteInternalError(w, "Failed to update person")
		return
	}

	updatedPerson, err := api.repo.GetByID(id)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	response := mappers.PersonDomainToResponse(updatedPerson)
	WriteSuccess(w, response)
}


// DeletePerson godoc
// @Summary Delete a person
// @Description Delete a person from the CRM system
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{id} [delete]
func (api *PersonAPI) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id, err := validators.ValidatePersonID(r.PathValue("id"))
	if err != nil {
		WriteBadRequest(w, err.Error())
		return
	}

	if err := api.repo.Delete(id); err != nil {
		WriteInternalError(w, "Failed to delete person")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
