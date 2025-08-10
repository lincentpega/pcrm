package repository

import (
    "database/sql"
    "errors"
    "fmt"

    "github.com/jmoiron/sqlx"
    "github.com/lincentpega/pcrm/internal/models"
)

type ConversationRepository struct {
    db *sqlx.DB
}

func NewConversationRepository(db *sqlx.DB) *ConversationRepository {
    return &ConversationRepository{db: db}
}

func (r *ConversationRepository) GetConversationTypes() ([]models.ConversationType, error) {
    var types []models.ConversationType
    query := `SELECT id, name, created_at FROM conversation_types ORDER BY name`
    if err := r.db.Select(&types, query); err != nil {
        return nil, fmt.Errorf("failed to get conversation types: %w", err)
    }
    return types, nil
}

func (r *ConversationRepository) GetByPersonID(personID int64) ([]models.Conversation, error) {
    var conversations []models.Conversation
    query := `
        SELECT c.id, c.person_id, c.conversation_type_id, c.initiator, c.notes, c.created_at, c.updated_at,
               ct.id as "conversation_type.id", ct.name as "conversation_type.name", ct.created_at as "conversation_type.created_at"
        FROM conversations c
        JOIN conversation_types ct ON c.conversation_type_id = ct.id
        WHERE c.person_id = $1
        ORDER BY c.created_at DESC
    `
    if err := r.db.Select(&conversations, query, personID); err != nil {
        return nil, fmt.Errorf("failed to get conversations for person %d: %w", personID, err)
    }
    return conversations, nil
}

func (r *ConversationRepository) GetByID(id int64) (*models.Conversation, error) {
    var conversation models.Conversation
    query := `
        SELECT c.id, c.person_id, c.conversation_type_id, c.initiator, c.notes, c.created_at, c.updated_at,
               ct.id as "conversation_type.id", ct.name as "conversation_type.name", ct.created_at as "conversation_type.created_at"
        FROM conversations c
        JOIN conversation_types ct ON c.conversation_type_id = ct.id
        WHERE c.id = $1
    `
    if err := r.db.Get(&conversation, query, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get conversation by id %d: %w", id, err)
    }
    return &conversation, nil
}

func (r *ConversationRepository) Create(conversation *models.Conversation) error {
    query := `
        INSERT INTO conversations (person_id, conversation_type_id, initiator, notes)
        VALUES (:person_id, :conversation_type_id, :initiator, :notes)
        RETURNING id, created_at, updated_at
    `
    rows, err := r.db.NamedQuery(query, conversation)
    if err != nil {
        return fmt.Errorf("failed to create conversation: %w", err)
    }
    defer rows.Close()
    if rows.Next() {
        if err := rows.Scan(&conversation.ID, &conversation.CreatedAt, &conversation.UpdatedAt); err != nil {
            return fmt.Errorf("failed to scan created conversation: %w", err)
        }
    }
    return nil
}

func (r *ConversationRepository) Update(conversation *models.Conversation) error {
    query := `
        UPDATE conversations
        SET conversation_type_id = :conversation_type_id, initiator = :initiator, notes = :notes, updated_at = NOW()
        WHERE id = :id
        RETURNING updated_at
    `
    rows, err := r.db.NamedQuery(query, conversation)
    if err != nil {
        return fmt.Errorf("failed to update conversation: %w", err)
    }
    defer rows.Close()
    if rows.Next() {
        if err := rows.Scan(&conversation.UpdatedAt); err != nil {
            return fmt.Errorf("failed to scan updated conversation: %w", err)
        }
    }
    return nil
}

func (r *ConversationRepository) Delete(id int64) error {
    query := `DELETE FROM conversations WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete conversation: %w", err)
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("conversation with id %d not found", id)
    }
    return nil
}

