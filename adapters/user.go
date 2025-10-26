package adapters

import (
	"accounter/domain/user"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	ctx  context.Context
	conn *sqlx.DB
}

func (r *UserRepository) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.ctx, time.Minute*5)
}

func NewUserRepository(ctx context.Context, conn *sqlx.DB) *UserRepository {
	return &UserRepository{ctx: ctx, conn: conn}
}

func (r *UserRepository) GetList() ([]user.User, error) {
	ctx, cancel := r.getContext()
	defer cancel()

	query := "SELECT * FROM users"
	result := make([]user.User, 0)

	err := r.conn.SelectContext(ctx, &result, query)

	return result, err
}

func (r *UserRepository) GetOne(id int64) (user.User, error) {
	return user.User{}, nil
}

func (r *UserRepository) Save(user *user.User) error {
	ctx, cancel := r.getContext()
	defer cancel()

	query := `
		INSERT INTO users
			(id, login, password, name, surname, patronymic, price_per_hour)
		VALUES (:id, :login, :password, :name, :surname, :patronymic, :price_per_hour)
			ON CONFLICT(id) DO UPDATE SET
				login=:login,
				password=:password,
				name=:name,
				surname=:surname,
				patronymic=:patronymic,
				price_per_hour=:price_per_hour
	`

	if res, err := r.conn.NamedExecContext(ctx, query, user); err != nil {
		return err
	} else if id, _ := res.LastInsertId(); id != 0 {
		user.ID = id
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	return nil
}
