package dto

import "time"

type ConversationRequest struct {
    ConversationTypeID int64  `json:"conversationTypeId"`
    Initiator          string `json:"initiator"`
    Notes              string `json:"notes"`
}

type ConversationResponse struct {
    ID                 int64                    `json:"id"`
    PersonID           int64                    `json:"personId"`
    Initiator          string                   `json:"initiator"`
    Notes              string                   `json:"notes"`
    CreatedAt          time.Time                `json:"createdAt"`
    UpdatedAt          time.Time                `json:"updatedAt"`
    ConversationType   ConversationTypeResponse `json:"conversationType"`
}

type ConversationTypeResponse struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}
