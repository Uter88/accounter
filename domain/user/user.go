package user

import (
	"accounter/tools"
)

type Users []User

// User model
type User struct {
	ID           int64   `db:"id,omitempty" json:"id"`
	Login        string  `db:"login" json:"login"`
	Password     string  `db:"password" json:"password"`
	Name         string  `db:"name" json:"name"`
	Surname      string  `db:"surname" json:"surname"`
	Patronymic   string  `db:"patronymic" json:"patronymic"`
	PricePerHour float32 `db:"price_per_hour" json:"price_per_hour"`
}

// IsValid check for User data is valid
func (u *User) IsValid(isAuth bool) bool {
	if err := tools.ValidEmail(u.Login); err != nil {
		return false
	}

	if isAuth {
		if tools.IsSomeEmpty(u.Login, u.Password) {
			return false
		}
	} else {
		if tools.IsSomeEmpty(u.Login, u.Password, u.Name, u.Surname, u.Patronymic) {
			return false
		}

		if tools.IsEmptyValue(u.PricePerHour) {
			return false
		}
	}

	return true
}

// Reset reset User data
func (u *User) Reset() {
	u = &User{}
}
