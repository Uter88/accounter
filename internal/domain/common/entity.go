package common

// Entity common entity with minimal identification params
type Entity interface {
	// Unique entity ID
	GetID() int64

	// Entity type
	GetType() string
}
