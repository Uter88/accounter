package postgres

import (
	"accounter/internal/domain/common"
	"accounter/internal/domain/user"
	"accounter/pkg/utils"
	"context"
	"database/sql"
	"errors"
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
func (r *userRepository) GetList(ctx context.Context, params common.RequestParams) (user.Users, error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	result := make(user.Users, 0)
	err := r.namedSelect(ctx, getUserListQuery, &result, params)

	return result, err
}

// GetOne get one User by id
func (r *userRepository) GetOne(ctx context.Context, id int64) (u user.User, err error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	db := r.client.GetExecutor(ctx)

	query := fmt.Sprintf(getUserQuery, "u.id = $1")
	err = db.GetContext(ctx, &u, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		err = user.ErrUserNotFound
	}

	return
}

// GetByCredentials get one User by login and/or password
func (r *userRepository) GetByCredentials(ctx context.Context, login, password string) (u user.User, err error) {
	ctx, cancel := r.getContext(ctx)
	defer cancel()

	params := utils.Data{
		"login":    login,
		"password": password,
	}

	cond := "u.login = :login"

	if password != "" {
		cond += " AND u.password = :password"
	}

	query := fmt.Sprintf(getUserQuery, cond)
	err = r.namedGet(ctx, query, &u, params)

	if errors.Is(err, sql.ErrNoRows) {
		err = user.ErrUserNotFound
	}

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
	getUserListQuery = `
		SELECT
			u.id,
			u.login,
			u.password,
			u.name,
			u.surname,
			u.patronymic,
			u.price_per_hour,
			COALESCE((
				SELECT SUM((CAST(t.work_end-t.work_begin as float) / 3600) * t.price_per_hour)
				FROM tasks t 
				WHERE t.user_id = u.id AND t.date BETWEEN :date_start AND :date_end
			), 0) as money_earned
		FROM users u
		GROUP BY u.id
		ORDER BY u.surname
	`
	getUserQuery = `
		SELECT
			u.id,
			u.login,
			u.password,
			u.name,
			u.surname,
			u.patronymic,
			u.price_per_hour
		FROM users u
		WHERE %s
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
