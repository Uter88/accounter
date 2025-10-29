package user

import (
	"accounter/tools"
)

type User struct {
	ID           int64   `db:"id,omitempty" json:"id"`
	Login        string  `db:"login" json:"login"`
	Password     string  `db:"password" json:"password"`
	Name         string  `db:"name" json:"name"`
	Surname      string  `db:"surname" json:"surname"`
	Patronymic   string  `db:"patronymic" json:"patronymic"`
	PricePerHour float64 `db:"price_per_hour" json:"price_per_hour"`
}

func (u *User) IsValid() bool {
	if tools.IsSomeEmpty(u.Login, u.Password, u.Name, u.Surname, u.Patronymic) {
		return false
	}

	if tools.IsEmptyValue(u.PricePerHour) {
		return false
	}

	return true
}

func (u *User) Reset() {
	u = &User{}
}

// Reflex
func (u *User) ResetField(field string) {
	switch field {
	case "name":
		u.Name = ""
	case "surname":
		u.Surname = ""
	case "patronymic":
		u.Patronymic = ""
	case "login":
		u.Login = ""
	case "password":
		u.Password = ""
	}
}
