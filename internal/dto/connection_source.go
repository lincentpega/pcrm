package dto

import "time"

type ConnectionSourceRequest struct {
	MeetingStory       *string    `json:"meetingStory,omitempty"`
	MeetingTimestamp   *time.Time `json:"meetingTimestamp,omitempty"`
	WasIntroduced      *bool      `json:"wasIntroduced,omitempty"`
	IntroducerPersonID *int64     `json:"introducerPersonId,omitempty"`
	IntroducerName     *string    `json:"introducerName,omitempty"`
}

type ConnectionSourceResponse struct {
	ID                 int64      `json:"id"`
	PersonID           int64      `json:"personId"`
	MeetingStory       *string    `json:"meetingStory,omitempty"`
	MeetingTimestamp   *time.Time `json:"meetingTimestamp,omitempty"`
	WasIntroduced      *bool      `json:"wasIntroduced,omitempty"`
	IntroducerPersonID *int64     `json:"introducerPersonId,omitempty"`
	IntroducerName     *string    `json:"introducerName,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}
