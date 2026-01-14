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
	ID         int64               `json:"id"`
	UserID     int64               `json:"user_id"`
	ObjectID   int64               `json:"object_id"`
	ObjectType string              `json:"object_type"`
	Date       time.Time           `json:"date"`
	Type       EventType           `json:"type"`
	Path       string              `json:"path"`
	Updates    []EventUpdateRecord `json:"updates"`
}

// EventUpdateRecord record with update field information
type EventUpdateRecord struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
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
