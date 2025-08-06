package models

import (
	"time"
)

type BirthDateInfo struct {
	ID                      int64      `db:"id"`
	PersonID                int64      `db:"person_id"`
	BirthYear               *int       `db:"birth_year"`
	BirthMonth              *int       `db:"birth_month"`
	BirthDay                *int       `db:"birth_day"`
	ApproximateAge          *int       `db:"approximate_age"`
	ApproximateAgeUpdatedAt *time.Time `db:"approximate_age_updated_at"`
	CreatedAt               time.Time  `db:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at"`
}