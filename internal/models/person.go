package models

import (
	"time"
)

type Person struct {
	ID         int64     `db:"id"`
	FirstName  string    `db:"first_name"`
	SecondName *string   `db:"second_name"`
	MiddleName *string   `db:"middle_name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
