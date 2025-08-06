package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lincentpega/pcrm/internal/models"
)

type ConnectionSourceRepository struct {
	db *sqlx.DB
}

func NewConnectionSourceRepository(db *sqlx.DB) *ConnectionSourceRepository {
	return &ConnectionSourceRepository{db: db}
}

func (r *ConnectionSourceRepository) GetByPersonID(personID int64) (*models.ConnectionSource, error) {
	var connectionSource models.ConnectionSource
	query := `
		SELECT id, person_id, meeting_story, meeting_timestamp, was_introduced,
		       introducer_person_id, introducer_name, created_at, updated_at
		FROM connection_sources
		WHERE person_id = $1
	`
	
	if err := r.db.Get(&connectionSource, query, personID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get connection source for person %d: %w", personID, err)
	}
	
	return &connectionSource, nil
}

func (r *ConnectionSourceRepository) Create(connectionSource *models.ConnectionSource) error {
	query := `
		INSERT INTO connection_sources (person_id, meeting_story, meeting_timestamp, was_introduced,
		                               introducer_person_id, introducer_name)
		VALUES (:person_id, :meeting_story, :meeting_timestamp, :was_introduced,
		        :introducer_person_id, :introducer_name)
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, connectionSource)
	if err != nil {
		return fmt.Errorf("failed to create connection source: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&connectionSource.ID, &connectionSource.CreatedAt, &connectionSource.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan created connection source: %w", err)
		}
	}
	
	return nil
}

func (r *ConnectionSourceRepository) Update(connectionSource *models.ConnectionSource) error {
	query := `
		UPDATE connection_sources 
		SET meeting_story = :meeting_story, meeting_timestamp = :meeting_timestamp,
		    was_introduced = :was_introduced, introducer_person_id = :introducer_person_id,
		    introducer_name = :introducer_name, updated_at = NOW()
		WHERE person_id = :person_id
		RETURNING id, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, connectionSource)
	if err != nil {
		return fmt.Errorf("failed to update connection source: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&connectionSource.ID, &connectionSource.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan updated connection source: %w", err)
		}
	}
	
	return nil
}

func (r *ConnectionSourceRepository) Upsert(connectionSource *models.ConnectionSource) error {
	query := `
		INSERT INTO connection_sources (person_id, meeting_story, meeting_timestamp, was_introduced,
		                               introducer_person_id, introducer_name)
		VALUES (:person_id, :meeting_story, :meeting_timestamp, :was_introduced,
		        :introducer_person_id, :introducer_name)
		ON CONFLICT (person_id) DO UPDATE SET
		    meeting_story = EXCLUDED.meeting_story,
		    meeting_timestamp = EXCLUDED.meeting_timestamp,
		    was_introduced = EXCLUDED.was_introduced,
		    introducer_person_id = EXCLUDED.introducer_person_id,
		    introducer_name = EXCLUDED.introducer_name,
		    updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, connectionSource)
	if err != nil {
		return fmt.Errorf("failed to upsert connection source: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&connectionSource.ID, &connectionSource.CreatedAt, &connectionSource.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan upserted connection source: %w", err)
		}
	}
	
	return nil
}

func (r *ConnectionSourceRepository) Delete(personID int64) error {
	query := `DELETE FROM connection_sources WHERE person_id = $1`
	
	result, err := r.db.Exec(query, personID)
	if err != nil {
		return fmt.Errorf("failed to delete connection source: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("connection source for person %d not found", personID)
	}
	
	return nil
}