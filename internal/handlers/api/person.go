package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lincentpega/pcrm/internal/models"
	"github.com/lincentpega/pcrm/internal/repository"
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

type PersonRequest struct {
	FirstName  string     `json:"firstName" validate:"required"`
	SecondName *string    `json:"secondName,omitempty"`
	MiddleName *string    `json:"middleName,omitempty"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
}

type PersonResponse struct {
	ID         int64      `json:"id" binding:"required"`
	FirstName  string     `json:"firstName" binding:"required"`
	SecondName *string    `json:"secondName,omitempty"`
	MiddleName *string    `json:"middleName,omitempty"`
	Birthdate  *time.Time `json:"birthdate,omitempty"`
	CreatedAt  time.Time  `json:"createdAt" binding:"required"`
	UpdatedAt  time.Time  `json:"updatedAt" binding:"required"`
}

type PersonWithContactsResponse struct {
	PersonResponse
	Contacts []ContactResponse `json:"contacts"`
}

func (req *PersonRequest) ToPerson() *models.Person {
	return &models.Person{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		MiddleName: req.MiddleName,
		Birthdate:  req.Birthdate,
	}
}

func PersonToResponse(person *models.Person) PersonResponse {
	return PersonResponse{
		ID:         person.ID,
		FirstName:  person.FirstName,
		SecondName: person.SecondName,
		MiddleName: person.MiddleName,
		Birthdate:  person.Birthdate,
		CreatedAt:  person.CreatedAt,
		UpdatedAt:  person.UpdatedAt,
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
// @Success 200 {object} PaginatedResponse[PersonResponse]
// @Failure 500 {object} ErrorResponse
// @Router /api/people [get]
func (api *PersonAPI) ListPeople(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

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

	response := make([]PersonResponse, len(people))
	for i, person := range people {
		response[i] = PersonToResponse(&person)
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
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
		return
	}

	person, err := api.repo.GetByID(id)
	if err != nil {
		WriteNotFound(w, "Person not found")
		return
	}

	response := PersonToResponse(person)
	WriteSuccess(w, response)
}

// GetPersonWithContacts godoc
// @Summary Get a person with all their contacts
// @Description Get detailed information about a person including all their contact information
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} PersonWithContactsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{id}/full [get]
func (api *PersonAPI) GetPersonWithContacts(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
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

	response := PersonWithContactsResponse{
		PersonResponse: PersonToResponse(person),
		Contacts:       make([]ContactResponse, len(contacts)),
	}

	for i, contact := range contacts {
		response.Contacts[i] = ContactToResponse(&contact)
	}

	WriteSuccess(w, response)
}

// CreatePerson godoc
// @Summary Create a new person
// @Description Create a new person in the CRM system
// @Tags people
// @Accept json
// @Produce json
// @Param person body PersonRequest true "Person data"
// @Success 201 {object} PersonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people [post]
func (api *PersonAPI) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if req.FirstName == "" {
		WriteBadRequest(w, "First name is required")
		return
	}

	person := req.ToPerson()
	if err := api.repo.Create(person); err != nil {
		WriteInternalError(w, "Failed to create person")
		return
	}

	response := PersonToResponse(person)
	WriteCreated(w, response)
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Update an existing person's information
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body PersonRequest true "Updated person data"
// @Success 200 {object} PersonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{id} [put]
func (api *PersonAPI) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
		return
	}

	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if req.FirstName == "" {
		WriteBadRequest(w, "First name is required")
		return
	}

	person := req.ToPerson()
	person.ID = id

	if err := api.repo.Update(person); err != nil {
		WriteInternalError(w, "Failed to update person")
		return
	}

	updatedPerson, err := api.repo.GetByID(id)
	if err != nil {
		WriteInternalError(w, "Failed to fetch updated person")
		return
	}

	response := PersonToResponse(updatedPerson)
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
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
		return
	}

	if err := api.repo.Delete(id); err != nil {
		WriteInternalError(w, "Failed to delete person")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}