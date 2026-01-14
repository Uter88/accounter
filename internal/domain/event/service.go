package event

import (
	"accounter/internal/domain/shared"
	"context"
)

// EventService event service
type EventService struct {
	repo       eventRepository
	comparator comporator

	pubslishers []EventPublisher
}

// NewEventService creates new EventService
func NewEventService(repo eventRepository, comparator comporator) *EventService {
	return &EventService{
		repo:       repo,
		comparator: comparator,
	}
}

// OnCreate trigger for creating events
func (es *EventService) OnCreate(ctx context.Context, userID int64, obj shared.Entity) error {
	objID := obj.GetID()
	objType := obj.GetType()

	event := NewEvent(userID, objID, objType, EventTypeCreate)

	return es.CreateEvent(ctx, event)
}

// OnUpdate trigger for updating events
func (es *EventService) OnUpdate(ctx context.Context, userID int64, old, new shared.Entity) error {
	objID := new.GetID()
	objType := new.GetType()
	updates, _ := es.comparator.Compare(old, new)

	event := NewEvent(userID, objID, objType, EventTypeUpdate, updates...)

	return es.CreateEvent(ctx, event)
}

// OnDelete trigger for deleting events
func (es *EventService) OnDelete(ctx context.Context, userID int64, obj shared.Entity) error {
	objID := obj.GetID()
	objType := obj.GetType()

	event := NewEvent(userID, objID, objType, EventTypeDelete)

	return es.CreateEvent(ctx, event)
}

// CreateEvent creates new Event
func (es *EventService) CreateEvent(ctx context.Context, event Event) error {
	if err := es.repo.Create(ctx, &event); err != nil {
		return err
	}

	for _, publisher := range es.pubslishers {
		publisher.PublishEvent(event)
	}

	return nil
}

// GetEventList get Events
func (es *EventService) GetEventList(ctx context.Context) (Events, error) {
	return es.repo.GetList(ctx)
}

// RegisterPublisher register new EventPublisher
func (es *EventService) RegisterPublisher(publishers EventPublisher) {
	es.pubslishers = append(es.pubslishers, publishers)
}
