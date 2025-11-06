package store

import (
	v1 "accounter/backend/server/handlers/v1"
	"accounter/domain/user"
	"accounter/tools"
)

type usersStore struct {
	baseStore
	users user.Users
}

func (s *usersStore) SaveUser(u user.User) error {
	api := s.getApi("users/save")

	resp, _, err := tools.MakeJSONRequest[v1.Response[user.User], any]("POST", api, tools.ToJSON(u))

	if err != nil {
		return err
	}

	s.UpdateUser(resp.Data)
	return nil
}

func (s *usersStore) RequestUsers() error {
	api := s.getApi("users/list")

	resp, _, err := tools.MakeJSONRequest[v1.Response[user.Users], any]("GET", api, nil)

	if err != nil {
		return err
	}

	s.users = resp.Data

	return nil
}

func (s *usersStore) GetUsers() user.Users {
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
