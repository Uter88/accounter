package user

// User port
type UserRepository interface {

	// Get list of User
	GetList() ([]User, error)

	// Get one User by id
	GetOne(id int64) (User, error)

	// Create/update User
	Save(*User) error

	// Delete one User by id
	Delete(id int64) error

	// Get one User by login
	GetByCredentials(login, password string) (User, error)
}
