package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lincentpega/pcrm/internal/models"
)

type ContactRepository struct {
	db *sqlx.DB
}

func NewContactRepository(db *sqlx.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) GetContactTypes() ([]models.ContactType, error) {
	var types []models.ContactType
	query := `SELECT id, name, created_at FROM contact_types ORDER BY name`
	
	if err := r.db.Select(&types, query); err != nil {
		return nil, fmt.Errorf("failed to get contact types: %w", err)
	}
	
	return types, nil
}

func (r *ContactRepository) GetByPersonID(personID int64) ([]models.Contact, error) {
	var contacts []models.Contact
	query := `
		SELECT c.id, c.person_id, c.contact_type_id, c.content, c.created_at, c.updated_at,
		       ct.id as "contact_type.id", ct.name as "contact_type.name", ct.created_at as "contact_type.created_at"
		FROM contacts c
		JOIN contact_types ct ON c.contact_type_id = ct.id
		WHERE c.person_id = $1
		ORDER BY ct.name, c.created_at DESC
	`
	
	if err := r.db.Select(&contacts, query, personID); err != nil {
		return nil, fmt.Errorf("failed to get contacts for person %d: %w", personID, err)
	}
	
	return contacts, nil
}

func (r *ContactRepository) GetByID(id int64) (*models.Contact, error) {
	var contact models.Contact
	query := `
		SELECT c.id, c.person_id, c.contact_type_id, c.content, c.created_at, c.updated_at,
		       ct.id as "contact_type.id", ct.name as "contact_type.name", ct.created_at as "contact_type.created_at"
		FROM contacts c
		JOIN contact_types ct ON c.contact_type_id = ct.id
		WHERE c.id = $1
	`
	
	if err := r.db.Get(&contact, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contact by id %d: %w", id, err)
	}
	
	return &contact, nil
}

func (r *ContactRepository) Create(contact *models.Contact) error {
	query := `
		INSERT INTO contacts (person_id, contact_type_id, content)
		VALUES (:person_id, :contact_type_id, :content)
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, contact)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan created contact: %w", err)
		}
	}
	
	return nil
}

func (r *ContactRepository) Update(contact *models.Contact) error {
	query := `
		UPDATE contacts 
		SET contact_type_id = :contact_type_id, content = :content, updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at
	`
	
	rows, err := r.db.NamedQuery(query, contact)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&contact.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan updated contact: %w", err)
		}
	}
	
	return nil
}

func (r *ContactRepository) Delete(id int64) error {
	query := `DELETE FROM contacts WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("contact with id %d not found", id)
	}
	
	return nil
}