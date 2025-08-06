package models

import (
	"time"
)

type ConnectionSource struct {
	ID                   int64      `db:"id"`
	PersonID             int64      `db:"person_id"`
	MeetingStory         *string    `db:"meeting_story"`
	MeetingTimestamp     *time.Time `db:"meeting_timestamp"`
	WasIntroduced        *bool      `db:"was_introduced"`
	IntroducerPersonID   *int64     `db:"introducer_person_id"`
	IntroducerName       *string    `db:"introducer_name"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}