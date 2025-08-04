package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/models"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) GetPaginated(page, limit int) ([]models.Person, error) {
	var people []models.Person
	offset := (page - 1) * limit
	
	query := `
		SELECT id, first_name, second_name, middle_name, 
		       birth_year, birth_month, birth_day, approximate_age, approximate_age_updated_at,
		       created_at, updated_at
		FROM people
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	if err := r.db.Select(&people, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get paginated people: %w", err)
	}
	
	return people, nil
}

func (r *PersonRepository) GetTotalCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM people`
	
	if err := r.db.Get(&count, query); err != nil {
		return 0, fmt.Errorf("failed to get people count: %w", err)
	}
	
	return count, nil
}

func (r *PersonRepository) GetByID(id int64) (*models.Person, error) {
	var person models.Person
	query := `
		SELECT id, first_name, second_name, middle_name,
		       birth_year, birth_month, birth_day, approximate_age, approximate_age_updated_at,
		       created_at, updated_at
		FROM people
		WHERE id = $1
	`
	
	if err := r.db.Get(&person, query, id); err != nil {
		return nil, fmt.Errorf("failed to get person by id %d: %w", id, err)
	}
	
	return &person, nil
}

func (r *PersonRepository) Create(person *models.Person) error {
	query := `
		INSERT INTO people (first_name, second_name, middle_name, 
		                   birth_year, birth_month, birth_day, approximate_age, approximate_age_updated_at)
		VALUES (:first_name, :second_name, :middle_name,
		        :birth_year, :birth_month, :birth_day, :approximate_age, :approximate_age_updated_at)
		RETURNING id, created_at, updated_at
	`
	
	rows, err := r.db.NamedQuery(query, person)
	if err != nil {
		return fmt.Errorf("failed to create person: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan created person: %w", err)
		}
	}
	
	return nil
}

func (r *PersonRepository) Update(person *models.Person) error {
	query := `
		UPDATE people 
		SET first_name = :first_name, second_name = :second_name, middle_name = :middle_name,
		    birth_year = :birth_year, birth_month = :birth_month, birth_day = :birth_day,
		    approximate_age = :approximate_age, approximate_age_updated_at = :approximate_age_updated_at,
		    updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at
	`
	
	rows, err := r.db.NamedQuery(query, person)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&person.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan updated person: %w", err)
		}
	}
	
	return nil
}

func (r *PersonRepository) Delete(id int64) error {
	query := `DELETE FROM people WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("person with id %d not found", id)
	}
	
	return nil
}


func (r *PersonRepository) UpdateBirthdate(personID int64, req *dto.PersonBirthdateRequest) error {
	var approximateAgeUpdatedAt *time.Time
	if req.ApproximateAge != nil {
		now := time.Now()
		approximateAgeUpdatedAt = &now
	}

	query := `
		UPDATE people 
		SET birth_year = $2, birth_month = $3, birth_day = $4,
		    approximate_age = $5, approximate_age_updated_at = $6,
		    updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.Exec(query, personID, req.BirthYear, req.BirthMonth, req.BirthDay, 
		                      req.ApproximateAge, approximateAgeUpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update birth date: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("person with id %d not found", personID)
	}
	
	return nil
}