package user

import "context"

// User port
type UserRepository interface {

	// Get list of User
	GetList(ctx context.Context) (Users, error)

	// Get one User by id
	GetOne(ctx context.Context, id int64) (User, error)

	// Create User
	Create(ctx context.Context, user *User) error

	// Update User
	Update(ctx context.Context, user *User) error

	// Delete one User by id
	Delete(ctx context.Context, id int64) error

	// Get one User by login
	GetByCredentials(ctx context.Context, login, password string) (User, error)

	// Execute operations in transaction
	WithTx(ctx context.Context, cb func(context.Context) error) error
}
