package adapter_sql

import (
	"accounter/internal/domain/user"
	"accounter/pkg/tools"
	"context"
	"fmt"
)

// User repository
type userRepository struct {
	*baseRepository
}

// NewUserRepository creates new userRepository
func NewUserRepository(client *SQLClient) *userRepository {
	return &userRepository{
		baseRepository: newBaseRepository(client),
	}
}

// GetList get list of User
func (r *userRepository) GetList(ctx context.Context) (user.Users, error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	result := make(user.Users, 0)
	query := fmt.Sprintf("%s ORDER BY login", getUserQuery)

	db := r.client.GetExecutor(ctx)
	err := db.SelectContext(ctx, &result, query)

	return result, err
}

// GetOne get one User by id
func (r *userRepository) GetOne(ctx context.Context, id int64) (u user.User, err error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	query := fmt.Sprintf("%s WHERE id = $1", getUserQuery)
	err = db.GetContext(ctx, &u, query, id)

	return
}

// GetByCredentials get one User by login and/or password
func (r *userRepository) GetByCredentials(ctx context.Context, login, password string) (u user.User, err error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	cond := "WHERE login = :login"

	if password != "" {
		cond += " AND password = :password"
	}

	query := fmt.Sprintf("%s %s", getUserQuery, cond)

	err = r.namedGet(ctx, query, &u, tools.Data{"login": login, "password": password})

	return
}

// Create creates new User
func (r *userRepository) Create(ctx context.Context, user *user.User) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	if res, err := db.NamedExecContext(ctx, insertUserQuery, user); err != nil {
		return err

	} else if id, _ := res.LastInsertId(); id != 0 {
		user.ID = id
	}

	return nil
}

// Update updates one User
func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	_, err := db.NamedExecContext(ctx, insertUserQuery, user)

	return err
}

// Delete deletes one User by id
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)
	_, err := db.ExecContext(ctx, deleteUserQuery, id)

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
		RETURNING id;
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
