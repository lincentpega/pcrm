package api

import (
    "encoding/json"
    "net/http"

    "github.com/lincentpega/pcrm/internal/dto"
    "github.com/lincentpega/pcrm/internal/mappers"
    "github.com/lincentpega/pcrm/internal/models"
    "github.com/lincentpega/pcrm/internal/repository"
    "github.com/lincentpega/pcrm/internal/validators"
)

type ConversationAPI struct {
    repo       *repository.ConversationRepository
    personRepo *repository.PersonRepository
}

func NewConversationAPI(repo *repository.ConversationRepository, personRepo *repository.PersonRepository) *ConversationAPI {
    return &ConversationAPI{repo: repo, personRepo: personRepo}
}

// ListConversationsByPerson godoc
// @Summary List conversations for a person
// @Description Get all conversations associated with a specific person
// @Tags conversations
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Success 200 {array} dto.ConversationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/conversations [get]
func (api *ConversationAPI) ListConversationsByPerson(w http.ResponseWriter, r *http.Request) {
    personID, err := validators.ValidatePersonID(r.PathValue("personId"))
    if err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    conversations, err := api.repo.GetByPersonID(personID)
    if err != nil {
        WriteInternalError(w, "Failed to fetch conversations")
        return
    }
    response := make([]dto.ConversationResponse, len(conversations))
    for i, c := range conversations {
        response[i] = mappers.ConversationDomainToResponse(&c)
    }
    WriteSuccess(w, response)
}

// GetConversation godoc
// @Summary Get a conversation by ID
// @Description Get detailed information about a specific conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Success 200 {object} dto.ConversationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/conversations/{id} [get]
func (api *ConversationAPI) GetConversation(w http.ResponseWriter, r *http.Request) {
    id, err := validators.ValidateConversationID(r.PathValue("id"))
    if err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    conversation, err := api.repo.GetByID(id)
    if err != nil {
        WriteNotFound(w, "Conversation not found")
        return
    }
    response := mappers.ConversationDomainToResponse(conversation)
    WriteSuccess(w, response)
}

// CreateConversation godoc
// @Summary Create a new conversation
// @Description Create a new conversation for a specific person
// @Tags conversations
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Param conversation body dto.ConversationRequest true "Conversation data"
// @Success 201 {object} dto.ConversationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/people/{personId}/conversations [post]
func (api *ConversationAPI) CreateConversation(w http.ResponseWriter, r *http.Request) {
    personID, err := validators.ValidatePersonID(r.PathValue("personId"))
    if err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    person, err := api.personRepo.GetByID(personID)
    if err != nil {
        WriteInternalError(w, "Failed to fetch person")
        return
    }
    if person == nil {
        WriteNotFound(w, "Person not found")
        return
    }
    var req dto.ConversationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteBadRequest(w, "Invalid JSON format")
        return
    }
    if err := validators.ValidateConversationRequest(&req); err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    conversation := mappers.ConversationRequestToDomain(personID, &req)
    if err := api.repo.Create(conversation); err != nil {
        WriteInternalError(w, "Failed to create conversation")
        return
    }
    created, err := api.repo.GetByID(conversation.ID)
    if err != nil || created == nil {
        WriteInternalError(w, "Failed to fetch created conversation")
        return
    }
    response := mappers.ConversationDomainToResponse(created)
    WriteCreated(w, response)
}

// UpdateConversation godoc
// @Summary Update a conversation
// @Description Update an existing conversation's information
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param conversation body dto.ConversationRequest true "Updated conversation data"
// @Success 200 {object} dto.ConversationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/conversations/{id} [put]
func (api *ConversationAPI) UpdateConversation(w http.ResponseWriter, r *http.Request) {
    id, err := validators.ValidateConversationID(r.PathValue("id"))
    if err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    var req dto.ConversationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteBadRequest(w, "Invalid JSON format")
        return
    }
    if err := validators.ValidateConversationRequest(&req); err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    existing, err := api.repo.GetByID(id)
    if err != nil || existing == nil {
        WriteNotFound(w, "Conversation not found")
        return
    }
    conversation := &models.Conversation{ID: id, PersonID: existing.PersonID, ConversationTypeID: req.ConversationTypeID, Initiator: req.Initiator, Notes: req.Notes}
    if err := api.repo.Update(conversation); err != nil {
        WriteInternalError(w, "Failed to update conversation")
        return
    }
    updated, err := api.repo.GetByID(id)
    if err != nil || updated == nil {
        WriteInternalError(w, "Failed to fetch updated conversation")
        return
    }
    response := mappers.ConversationDomainToResponse(updated)
    WriteSuccess(w, response)
}

// DeleteConversation godoc
// @Summary Delete a conversation
// @Description Delete a conversation from the system
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/conversations/{id} [delete]
func (api *ConversationAPI) DeleteConversation(w http.ResponseWriter, r *http.Request) {
    id, err := validators.ValidateConversationID(r.PathValue("id"))
    if err != nil {
        WriteBadRequest(w, err.Error())
        return
    }
    if err := api.repo.Delete(id); err != nil {
        WriteInternalError(w, "Failed to delete conversation")
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// ListConversationTypes godoc
// @Summary List all conversation types
// @Description Get all available conversation types (phone call, video call, etc.)
// @Tags conversation-types
// @Accept json
// @Produce json
// @Success 200 {array} dto.ConversationTypeResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/conversation-types [get]
func (api *ConversationAPI) ListConversationTypes(w http.ResponseWriter, r *http.Request) {
    types, err := api.repo.GetConversationTypes()
    if err != nil {
        WriteInternalError(w, "Failed to fetch conversation types")
        return
    }
    response := make([]dto.ConversationTypeResponse, len(types))
    for i, t := range types {
        response[i] = dto.ConversationTypeResponse{ID: t.ID, Name: t.Name, CreatedAt: t.CreatedAt}
    }
    WriteSuccess(w, response)
}

