package models

import (
	"time"
)

type ContactType struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Contact struct {
	ID            int64       `json:"id" db:"id"`
	PersonID      int64       `json:"personId" db:"person_id"`
	ContactTypeID int64       `json:"contactTypeId" db:"contact_type_id"`
	Content       string      `json:"content" db:"content"`
	CreatedAt     time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time   `json:"updatedAt" db:"updated_at"`
	ContactType   ContactType `json:"contactType" db:"contact_type"`
}

type ContactForm struct {
	ContactTypeID int64  `json:"contactTypeId"`
	Content       string `json:"content"`
}

func (f *ContactForm) ToContact(personID int64) *Contact {
	return &Contact{
		PersonID:      personID,
		ContactTypeID: f.ContactTypeID,
		Content:       f.Content,
	}
}

func (c *Contact) ToForm() *ContactForm {
	return &ContactForm{
		ContactTypeID: c.ContactTypeID,
		Content:       c.Content,
	}
}