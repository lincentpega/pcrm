package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Person struct {
	ID         int64          `json:"id" db:"id"`
	FirstName  string         `json:"firstName" db:"first_name"`
	SecondName sql.NullString `json:"secondName" db:"second_name"`
	MiddleName sql.NullString `json:"middleName" db:"middle_name"`
	Birthdate  sql.NullTime   `json:"birthdate" db:"birthdate"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time      `json:"updatedAt" db:"updated_at"`
}

func (p Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	aux := struct {
		SecondName *string `json:"secondName,omitempty"`
		MiddleName *string `json:"middleName,omitempty"`
		Birthdate  *string `json:"birthdate,omitempty"`
		Alias
	}{
		Alias: Alias(p),
	}

	if p.SecondName.Valid {
		aux.SecondName = &p.SecondName.String
	}
	if p.MiddleName.Valid {
		aux.MiddleName = &p.MiddleName.String
	}
	if p.Birthdate.Valid {
		formatted := p.Birthdate.Time.Format("2006-01-02")
		aux.Birthdate = &formatted
	}

	return json.Marshal(aux)
}

type PersonForm struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	MiddleName string `json:"middleName"`
	Birthdate  string `json:"birthdate"`
}

func (p *PersonForm) ToPerson() *Person {
	person := &Person{
		FirstName: p.FirstName,
	}

	if p.SecondName != "" {
		person.SecondName = sql.NullString{String: p.SecondName, Valid: true}
	}

	if p.MiddleName != "" {
		person.MiddleName = sql.NullString{String: p.MiddleName, Valid: true}
	}

	if p.Birthdate != "" {
		if birthdate, err := time.Parse("2006-01-02", p.Birthdate); err == nil {
			person.Birthdate = sql.NullTime{Time: birthdate, Valid: true}
		}
	}

	return person
}

func (p *Person) ToForm() *PersonForm {
	form := &PersonForm{
		FirstName: p.FirstName,
	}

	if p.SecondName.Valid {
		form.SecondName = p.SecondName.String
	}

	if p.MiddleName.Valid {
		form.MiddleName = p.MiddleName.String
	}

	if p.Birthdate.Valid {
		form.Birthdate = p.Birthdate.Time.Format("2006-01-02")
	}

	return form
}