package models

import "accounter/domain/user"

type LoginForm struct {
	user.User
	IsRemember bool
	IsAccept   bool
}

func (f *LoginForm) Validate(isAuth bool) bool {
	if !f.IsValid(isAuth) {
		return false
	}

	if !isAuth && !f.IsAccept {
		return false
	}

	return true
}
