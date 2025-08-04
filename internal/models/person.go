package models

import (
	"time"
)

type Person struct {
	ID                      int64      `db:"id"`
	FirstName               string     `db:"first_name"`
	SecondName              *string    `db:"second_name"`
	MiddleName              *string    `db:"middle_name"`
	BirthYear               *int       `db:"birth_year"`
	BirthMonth              *int       `db:"birth_month"`
	BirthDay                *int       `db:"birth_day"`
	ApproximateAge          *int       `db:"approximate_age"`
	ApproximateAgeUpdatedAt *time.Time `db:"approximate_age_updated_at"`
	CreatedAt               time.Time  `db:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at"`
}
