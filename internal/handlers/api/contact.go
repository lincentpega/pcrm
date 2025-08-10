package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lincentpega/pcrm/internal/models"
	"github.com/lincentpega/pcrm/internal/repository"
)

type ContactAPI struct {
	contactRepo *repository.ContactRepository
}

func NewContactAPI(contactRepo *repository.ContactRepository) *ContactAPI {
	return &ContactAPI{
		contactRepo: contactRepo,
	}
}

type ContactRequest struct {
	ContactTypeID int64  `json:"contactTypeId" validate:"required"`
	Content       string `json:"content" validate:"required"`
}

type ContactResponse struct {
    ID            int64                `json:"id"`
    PersonID      int64                `json:"personId"`
    Content       string               `json:"content"`
    CreatedAt     time.Time            `json:"createdAt"`
    UpdatedAt     time.Time            `json:"updatedAt"`
    ContactType   ContactTypeResponse  `json:"contactType"`
}

type ContactTypeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (req *ContactRequest) ToContact(personID int64) *models.Contact {
	return &models.Contact{
		PersonID:      personID,
		ContactTypeID: req.ContactTypeID,
		Content:       req.Content,
	}
}

func ContactToResponse(contact *models.Contact) ContactResponse {
    return ContactResponse{
        ID:            contact.ID,
        PersonID:      contact.PersonID,
        Content:       contact.Content,
        CreatedAt:     contact.CreatedAt,
        UpdatedAt:     contact.UpdatedAt,
        ContactType: ContactTypeResponse{
            ID:        contact.ContactType.ID,
			Name:      contact.ContactType.Name,
			CreatedAt: contact.ContactType.CreatedAt,
		},
	}
}

func ContactTypeToResponse(contactType *models.ContactType) ContactTypeResponse {
	return ContactTypeResponse{
		ID:        contactType.ID,
		Name:      contactType.Name,
		CreatedAt: contactType.CreatedAt,
	}
}

// ListContactsByPerson godoc
// @Summary List contacts for a person
// @Description Get all contacts associated with a specific person
// @Tags contacts
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 200 {array} ContactResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/contacts [get]
func (api *ContactAPI) ListContactsByPerson(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
		return
	}

	contacts, err := api.contactRepo.GetByPersonID(personID)
	if err != nil {
		WriteInternalError(w, "Failed to fetch contacts")
		return
	}

	response := make([]ContactResponse, len(contacts))
	for i, contact := range contacts {
		response[i] = ContactToResponse(&contact)
	}

	WriteSuccess(w, response)
}

// GetContact godoc
// @Summary Get a contact by ID
// @Description Get detailed information about a specific contact
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 200 {object} ContactResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/contacts/{id} [get]
func (api *ContactAPI) GetContact(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        WriteBadRequest(w, "Invalid contact ID")
        return
    }

    contact, err := api.contactRepo.GetByID(id)
    if err != nil {
        WriteInternalError(w, "Failed to fetch contact")
        return
    }
    if contact == nil {
        WriteNotFound(w, "Contact not found")
        return
    }

    response := ContactToResponse(contact)
    WriteSuccess(w, response)
}

// CreateContact godoc
// @Summary Create a new contact
// @Description Create a new contact for a specific person
// @Tags contacts
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Param contact body ContactRequest true "Contact data"
// @Success 201 {object} ContactResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/contacts [post]
func (api *ContactAPI) CreateContact(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid person ID")
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if req.ContactTypeID == 0 {
		WriteBadRequest(w, "Contact type ID is required")
		return
	}

	if req.Content == "" {
		WriteBadRequest(w, "Content is required")
		return
	}

	contact := req.ToContact(personID)
	if err := api.contactRepo.Create(contact); err != nil {
		WriteInternalError(w, "Failed to create contact")
		return
	}

	createdContact, err := api.contactRepo.GetByID(contact.ID)
	if err != nil {
		WriteInternalError(w, "Failed to fetch created contact")
		return
	}

	response := ContactToResponse(createdContact)
	WriteCreated(w, response)
}

// UpdateContact godoc
// @Summary Update a contact
// @Description Update an existing contact's information
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Param contact body ContactRequest true "Updated contact data"
// @Success 200 {object} ContactResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/contacts/{id} [put]
func (api *ContactAPI) UpdateContact(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid contact ID")
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid JSON format")
		return
	}

	if req.ContactTypeID == 0 {
		WriteBadRequest(w, "Contact type ID is required")
		return
	}

	if req.Content == "" {
		WriteBadRequest(w, "Content is required")
		return
	}

    existingContact, err := api.contactRepo.GetByID(id)
    if err != nil {
        WriteInternalError(w, "Failed to fetch contact")
        return
    }
    if existingContact == nil {
        WriteNotFound(w, "Contact not found")
        return
    }

	contact := &models.Contact{
		ID:            id,
		PersonID:      existingContact.PersonID,
		ContactTypeID: req.ContactTypeID,
		Content:       req.Content,
	}

	if err := api.contactRepo.Update(contact); err != nil {
		WriteInternalError(w, "Failed to update contact")
		return
	}

	updatedContact, err := api.contactRepo.GetByID(id)
	if err != nil {
		WriteInternalError(w, "Failed to fetch updated contact")
		return
	}

	response := ContactToResponse(updatedContact)
	WriteSuccess(w, response)
}

// DeleteContact godoc
// @Summary Delete a contact
// @Description Delete a contact from the system
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/contacts/{id} [delete]
func (api *ContactAPI) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteBadRequest(w, "Invalid contact ID")
		return
	}

	if err := api.contactRepo.Delete(id); err != nil {
		WriteInternalError(w, "Failed to delete contact")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListContactTypes godoc
// @Summary List all contact types
// @Description Get all available contact types (email, phone, etc.)
// @Tags contact-types
// @Accept json
// @Produce json
// @Success 200 {array} ContactTypeResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/contact-types [get]
func (api *ContactAPI) ListContactTypes(w http.ResponseWriter, r *http.Request) {
	contactTypes, err := api.contactRepo.GetContactTypes()
	if err != nil {
		WriteInternalError(w, "Failed to fetch contact types")
		return
	}

	response := make([]ContactTypeResponse, len(contactTypes))
	for i, contactType := range contactTypes {
		response[i] = ContactTypeToResponse(&contactType)
	}

	WriteSuccess(w, response)
}
