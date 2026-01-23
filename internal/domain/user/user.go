package user

import (
	"accounter/pkg/utils"
	"fmt"
)

const UserType = "user"

// Users
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
	MoneyEarned  float32 `db:"money_earned" json:"money_earned"`
}

// GetLabel return short string User representation
func (u User) GetLabel() string {
	return fmt.Sprintf("%s %.1s. %.1s.", u.Surname, u.Name, u.Patronymic)
}

// IsValid check for User data is valid
func (u User) IsValid(isAuth bool) bool {
	if err := utils.ValidEmail(u.Login); err != nil {
		return false
	}

	if isAuth {
		if utils.IsSomeEmpty(u.Login, u.Password) {
			return false
		}
	} else {
		if utils.IsSomeEmpty(u.Login, u.Password, u.Name, u.Surname, u.Patronymic) {
			return false
		}

		if utils.IsEmptyValue(u.PricePerHour) {
			return false
		}
	}

	return true
}

// MoneyEarned return total money earned value
func (u User) GetMoneyEarned() float32 {
	return utils.ToFixed(u.MoneyEarned, 2)
}

// FormatMoneyEarned return total money earned value as string
func (u User) FormatMoneyEarned() string {
	return utils.FormatMoney(u.GetMoneyEarned())
}

// GetID get id
func (u User) GetID() int64 {
	return u.ID
}

// GetID get type of User entity
func (u User) GetType() string {
	return UserType
}
