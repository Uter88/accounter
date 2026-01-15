package event

import (
	"fmt"
	"time"
)

// Event types
type EventType = string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

// Events list of Event
type Events []Event

// Event for domain operations
type Event struct {
	ID         int64
	UserID     int64
	ObjectID   int64
	ObjectType string
	Date       time.Time
	Type       EventType
	Path       string
	Updates    []EventUpdateRecord
}

// EventUpdateRecord record with update field information
type EventUpdateRecord struct {
	Field    string
	OldValue any
	NewValue any
}

// NewEvent creates new Event
func NewEvent(userID, objectID int64, objectType string, eventType EventType, updates ...EventUpdateRecord) Event {
	return Event{
		UserID:     userID,
		ObjectID:   objectID,
		ObjectType: objectType,
		Type:       eventType,
		Path:       fmt.Sprintf("%s.%s.%d", objectType, eventType, objectID),
		Date:       time.Now(),
		Updates:    updates,
	}
}
