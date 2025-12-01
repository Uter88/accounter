package adapter_sql

import (
	"accounter/domain/user"
	"accounter/pkg/tools"
	"context"
	"fmt"
)

// User repository
type userRepository struct {
	baseRepository
}

// Creates new userRepository
func NewUserRepository(ctx context.Context, client SQLClient) *userRepository {
	return &userRepository{
		baseRepository: newBaseRepository(ctx, client),
	}
}

// Get list of User
func (r *userRepository) GetList() ([]user.User, error) {
	ctx, cancel := r.getContext()
	defer cancel()

	result := make([]user.User, 0)
	query := fmt.Sprintf("%s ORDER BY login", getUserQuery)
	err := r.db().SelectContext(ctx, &result, query)

	return result, err
}

// Get one User by id
func (r *userRepository) GetOne(id int64) (u user.User, err error) {
	ctx, cancel := r.getContext()
	defer cancel()

	query := fmt.Sprintf("%s WHERE id = $1", getUserQuery)
	err = r.db().GetContext(ctx, &u, query, id)

	return
}

// Get one User by login
func (r *userRepository) GetByCredentials(login, password string) (u user.User, err error) {
	ctx, cancel := r.getContext()
	defer cancel()

	cond := "WHERE login = :login"

	if password != "" {
		cond += " AND password = :password"
	}

	query := fmt.Sprintf("%s %s", getUserQuery, cond)

	err = r.namedGet(ctx, query, &u, tools.Data{"login": login, "password": password})

	return
}

// Save User
func (r *userRepository) Insert(user *user.User) error {
	ctx, cancel := r.getContext()
	defer cancel()

	if res, err := r.db().NamedExecContext(ctx, insertUserQuery, user); err != nil {
		return err

	} else if id, _ := res.LastInsertId(); id != 0 {
		user.ID = id
	}

	return nil
}

// Update one User
func (r *userRepository) Update(user *user.User) error {
	ctx, cancel := r.getContext()
	defer cancel()

	_, err := r.db().NamedExecContext(ctx, insertUserQuery, user)

	return err
}

// Delete one User by id
func (r *userRepository) Delete(id int64) error {
	ctx, cancel := r.getContext()
	defer cancel()

	_, err := r.db().ExecContext(ctx, deleteUserQuery, id)

	return err
}

// User queries
const (
	getUserQuery = `
		SELECT id, login, password, name, surname, patronymic, price_per_hour FROM users
	`
	deleteUserQuery = `DELETE FROM users WHERE id = $1`
	insertUserQuery = `
		INSERT INTO users (login, password, name, surname, patronymic, price_per_hour)
		VALUES (:id, :login, :password, :name, :surname, :patronymic, :price_per_hour)
	`
	updateUserQuery = `
		UPDATE users SET
			login=:login,
			password=:password,
			name=:name,
			surname=:surname,
			patronymic=:patronymic,
			price_per_hour=:price_per_hour
		WHERE id = :id
	`
)
