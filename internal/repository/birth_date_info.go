package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lincentpega/pcrm/internal/models"
)

type BirthDateInfoRepository struct {
	db *sqlx.DB
}

func NewBirthDateInfoRepository(db *sqlx.DB) *BirthDateInfoRepository {
	return &BirthDateInfoRepository{db: db}
}

func (r *BirthDateInfoRepository) GetByPersonID(personID int64) (*models.BirthDateInfo, error) {
	var birthDateInfo models.BirthDateInfo
	query := `
		SELECT id, person_id, birth_year, birth_month, birth_day,
		       approximate_age, approximate_age_updated_at, created_at, updated_at
		FROM birth_date_info
		WHERE person_id = $1
	`

	if err := r.db.Get(&birthDateInfo, query, personID); err != nil {
		return nil, fmt.Errorf("failed to get birth date info for person %d: %w", personID, err)
	}

	return &birthDateInfo, nil
}

func (r *BirthDateInfoRepository) Create(birthDateInfo *models.BirthDateInfo) error {
	query := `
		INSERT INTO birth_date_info (person_id, birth_year, birth_month, birth_day,
		                            approximate_age, approximate_age_updated_at)
		VALUES (:person_id, :birth_year, :birth_month, :birth_day,
		        :approximate_age, :approximate_age_updated_at)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQuery(query, birthDateInfo)
	if err != nil {
		return fmt.Errorf("failed to create birth date info: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&birthDateInfo.ID, &birthDateInfo.CreatedAt, &birthDateInfo.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan created birth date info: %w", err)
		}
	}

	return nil
}

func (r *BirthDateInfoRepository) Update(birthDateInfo *models.BirthDateInfo) error {
	query := `
		UPDATE birth_date_info 
		SET birth_year = :birth_year, birth_month = :birth_month, birth_day = :birth_day,
		    approximate_age = :approximate_age, approximate_age_updated_at = :approximate_age_updated_at,
		    updated_at = NOW()
		WHERE person_id = :person_id
		RETURNING id, updated_at
	`

	rows, err := r.db.NamedQuery(query, birthDateInfo)
	if err != nil {
		return fmt.Errorf("failed to update birth date info: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&birthDateInfo.ID, &birthDateInfo.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan updated birth date info: %w", err)
		}
	}

	return nil
}

func (r *BirthDateInfoRepository) Upsert(birthDateInfo *models.BirthDateInfo) error {
	query := `
		INSERT INTO birth_date_info (person_id, birth_year, birth_month, birth_day,
		                            approximate_age, approximate_age_updated_at)
		VALUES (:person_id, :birth_year, :birth_month, :birth_day,
		        :approximate_age, :approximate_age_updated_at)
		ON CONFLICT (person_id) DO UPDATE SET
		    birth_year = EXCLUDED.birth_year,
		    birth_month = EXCLUDED.birth_month,
		    birth_day = EXCLUDED.birth_day,
		    approximate_age = EXCLUDED.approximate_age,
		    approximate_age_updated_at = EXCLUDED.approximate_age_updated_at,
		    updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQuery(query, birthDateInfo)
	if err != nil {
		return fmt.Errorf("failed to upsert birth date info: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&birthDateInfo.ID, &birthDateInfo.CreatedAt, &birthDateInfo.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan upserted birth date info: %w", err)
		}
	}

	return nil
}

func (r *BirthDateInfoRepository) Delete(personID int64) error {
	query := `DELETE FROM birth_date_info WHERE person_id = $1`

	result, err := r.db.Exec(query, personID)
	if err != nil {
		return fmt.Errorf("failed to delete birth date info: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("birth date info for person %d not found", personID)
	}

	return nil
}
