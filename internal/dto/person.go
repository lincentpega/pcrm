package dto

import "time"

type PersonUpsertRequest struct {
	FirstName  string  `json:"firstName"`
	SecondName *string `json:"secondName,omitempty"`
	MiddleName *string `json:"middleName,omitempty"`
}

type PersonInfoResponse struct {
	ID         int64     `json:"id" binding:"required"`
	FirstName  string    `json:"firstName" binding:"required"`
	SecondName *string   `json:"secondName,omitempty"`
	MiddleName *string   `json:"middleName,omitempty"`
	CreatedAt  time.Time `json:"createdAt" binding:"required"`
	UpdatedAt  time.Time `json:"updatedAt" binding:"required"`
}

type PersonWithContactsResponse struct {
	PersonInfoResponse
	Contacts []ContactResponse `json:"contacts"`
}

type ContactResponse struct {
	ID          int64               `json:"id"`
	PersonID    int64               `json:"personId"`
	Content     string              `json:"content"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	ContactType ContactTypeResponse `json:"contactType"`
}

type ContactTypeResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
