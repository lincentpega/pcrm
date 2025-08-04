package models

import (
	"time"
)

type Person struct {
	ID         int64      `db:"id"`
	FirstName  string     `db:"first_name"`
	SecondName *string    `db:"second_name"`
	MiddleName *string    `db:"middle_name"`
	Birthdate  *time.Time `db:"birthdate"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
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
		person.SecondName = &p.SecondName
	}

	if p.MiddleName != "" {
		person.MiddleName = &p.MiddleName
	}

	if p.Birthdate != "" {
		if birthdate, err := time.Parse("2006-01-02", p.Birthdate); err == nil {
			person.Birthdate = &birthdate
		}
	}

	return person
}

func (p *Person) ToForm() *PersonForm {
	form := &PersonForm{
		FirstName: p.FirstName,
	}

	if p.SecondName != nil {
		form.SecondName = *p.SecondName
	}

	if p.MiddleName != nil {
		form.MiddleName = *p.MiddleName
	}

	if p.Birthdate != nil {
		form.Birthdate = p.Birthdate.Format("2006-01-02")
	}

	return form
}
