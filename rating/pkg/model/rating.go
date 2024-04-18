package model

type RecordID string
type RecordType string
type UserID string
type RatingValue int
type RatingEventType string

const (
	RecordTypeMovie = RecordType("movie")
)

const (
	RatingEventTypePut = "put"
	RatingEventTypeDelete = "delete"
)

type Rating struct {
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

// RatingEvent defines an event containing rating information.
type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}
