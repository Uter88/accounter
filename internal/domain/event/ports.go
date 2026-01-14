package event

import "context"

// Repository of Event
type eventRepository interface {
	GetList(ctx context.Context) (Events, error)
	Create(ctx context.Context, e *Event) error
}

// EventPublisher Event publisher
type EventPublisher interface {
	Name() string
	PublishEvent(event Event) error
}

// EventSubscriber Event subscriber
type EventSubscriber interface {
	Name() string
	SubscribeEvent(event Event) error
}

// comporator of changes between two same object
type comporator interface {
	Compare(a, b any) (result []EventUpdateRecord, ok bool)
}
