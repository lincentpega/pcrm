package dto

import "time"

type BirthDateInfoRequest struct {
	BirthYear      *int `json:"birthYear,omitempty"`
	BirthMonth     *int `json:"birthMonth,omitempty"`
	BirthDay       *int `json:"birthDay,omitempty"`
	ApproximateAge *int `json:"approximateAge,omitempty"`
}

type BirthDateInfoResponse struct {
	ID                      int64      `json:"id"`
	PersonID                int64      `json:"personId"`
	BirthYear               *int       `json:"birthYear,omitempty"`
	BirthMonth              *int       `json:"birthMonth,omitempty"`
	BirthDay                *int       `json:"birthDay,omitempty"`
	ApproximateAge          *int       `json:"approximateAge,omitempty"`
	ApproximateAgeUpdatedAt *time.Time `json:"approximateAgeUpdatedAt,omitempty"`
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               time.Time  `json:"updatedAt"`
}