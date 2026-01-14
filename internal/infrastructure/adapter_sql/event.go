package adapter_sql

import (
	"accounter/internal/domain/event"
	"accounter/pkg/tools"
	"context"
	"time"
)

// Event repository
type eventRepository struct {
	*baseRepository
}

// NewEventRepository creates new eventRepository
func NewEventRepository(client *SQLClient) *eventRepository {
	return &eventRepository{
		baseRepository: newBaseRepository(client),
	}
}

// GetList get list of Events
func (er *eventRepository) GetList(ctx context.Context) (event.Events, error) {
	ctx, cancel := er.getContext(ctx)
	defer cancel()

	db := er.client.GetExecutor(ctx)
	result := make(event.Events, 0)

	rows, err := db.QueryxContext(ctx, getEventQuery)

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var e eventDTO

		if err := rows.StructScan(&e); err != nil {
			return result, err
		}

		result = append(result, event.Event{
			UserID:     e.UserID,
			ObjectID:   e.ObjectID,
			ObjectType: e.ObjectType,
			Type:       e.Type,
			Date:       e.Date,
			Path:       e.Path,
			Updates:    tools.FromJSON[[]event.EventUpdateRecord](e.Updates),
		})
	}

	return nil, nil
}

// Create creates new Event
func (er *eventRepository) Create(ctx context.Context, e *event.Event) error {
	ctx, cancel := er.getContext(ctx)
	defer cancel()

	db := er.client.GetExecutor(ctx)

	data := eventDTO{
		UserID:     e.UserID,
		ObjectID:   e.ObjectID,
		ObjectType: e.ObjectType,
		Type:       e.Type,
		Date:       e.Date,
		Path:       e.Path,
		Updates:    tools.ToJSON(e.Updates).Bytes(),
	}

	if res, err := db.NamedExecContext(ctx, insertEventQuery, data); err != nil {
		return err

	} else if id, _ := res.LastInsertId(); id != 0 {
		e.ID = id
	}

	return nil
}

// Event data transfer struct
type eventDTO struct {
	ID         int64     `db:"id"`
	UserID     int64     `db:"user_id"`
	ObjectID   int64     `db:"object_id"`
	ObjectType string    `db:"object_type"`
	Type       string    `db:"type"`
	Path       string    `db:"path"`
	Date       time.Time `db:"date"`
	Updates    []byte    `db:"updates"`
}

// Event queries
const (
	insertEventQuery = `
		INSERT INTO events
			(user_id, object_id, object_type, date, type, path, updates)
		VALUES (:user_id, :object_id, :object_id, :date, :type, :path, :updates)
		RETURNING id;
	`

	getEventQuery = `
		SELECT user_id, object_id, object_type, date, type, path, updates
		FROM events
		ORDER BY date
	`
)
