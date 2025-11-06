package store

import (
	"accounter/domain/user"
	"accounter/tools"
	"net/http"
)

type usersStore struct {
	baseStore
	users []user.User
}

func (s *usersStore) SaveUser(u user.User) (user.User, error) {
	request := newRequest[user.User](s.baseStore).
		Path("users/save").
		Method(http.MethodPost).
		Data(tools.ToJSON(u))

	resp, _, err := request.Do()

	if err != nil {
		return u, err
	}

	result := resp.Data

	s.UpdateUser(result)
	return result, nil
}

func (s *usersStore) RequestUsers() error {
	request := newRequest[[]user.User](s.baseStore).
		Path("users/list")

	resp, _, err := request.Do()

	if err != nil {
		return err
	}

	s.users = resp.Data

	return nil
}

func (s *usersStore) CheckUniqueLogin(login string) (bool, error) {
	request := newRequest[bool](s.baseStore).
		Path("users/is_exists").
		Param("login", login)

	resp, _, err := request.Do()

	if err != nil {
		return false, err
	}

	return resp.Data, nil
}

func (s *usersStore) GetUsers() []user.User {
	return s.users
}

func (s *usersStore) UpdateUser(u user.User) {
	for i := range s.users {
		if s.users[i].ID == u.ID {
			s.users[i] = u
			return
		}
	}

	s.users = append(s.users, u)
}
