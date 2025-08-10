package mappers

import (
    "github.com/lincentpega/pcrm/internal/dto"
    "github.com/lincentpega/pcrm/internal/models"
)

func ConversationRequestToDomain(personID int64, req *dto.ConversationRequest) *models.Conversation {
    return &models.Conversation{
        PersonID:           personID,
        ConversationTypeID: req.ConversationTypeID,
        Initiator:          req.Initiator,
        Notes:              req.Notes,
    }
}

func ConversationDomainToResponse(conversation *models.Conversation) dto.ConversationResponse {
    return dto.ConversationResponse{
        ID:                 conversation.ID,
        PersonID:           conversation.PersonID,
        Initiator:          conversation.Initiator,
        Notes:              conversation.Notes,
        CreatedAt:          conversation.CreatedAt,
        UpdatedAt:          conversation.UpdatedAt,
        ConversationType: dto.ConversationTypeResponse{
            ID:        conversation.ConversationType.ID,
            Name:      conversation.ConversationType.Name,
            CreatedAt: conversation.ConversationType.CreatedAt,
        },
    }
}
