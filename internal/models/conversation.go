package models

import (
    "time"
)

type ConversationType struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Conversation struct {
    ID                  int64             `json:"id" db:"id"`
    PersonID            int64             `json:"personId" db:"person_id"`
    ConversationTypeID  int64             `json:"conversationTypeId" db:"conversation_type_id"`
    Initiator           string            `json:"initiator" db:"initiator"`
    Notes               string            `json:"notes" db:"notes"`
    CreatedAt           time.Time         `json:"createdAt" db:"created_at"`
    UpdatedAt           time.Time         `json:"updatedAt" db:"updated_at"`
    ConversationType    ConversationType  `json:"conversationType" db:"conversation_type"`
}

