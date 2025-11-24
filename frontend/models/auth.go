package models

import "accounter/domain/user"

type LoginForm struct {
	*user.User
	IsRemember bool
	IsAccept   bool
}

func NewLoginForm() LoginForm {
	return LoginForm{
		User: &user.User{},
	}
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

func (f *LoginForm) Reset() {
	f.User = &user.User{}
	f.IsRemember = false
	f.IsAccept = false
}
