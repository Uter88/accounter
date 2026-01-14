package store

import (
	"accounter/internal/domain/user"
	"accounter/pkg/utils"
	"fmt"
	"net/http"
)

type usersStore struct {
	*baseStore
	users   user.Users
	loading bool
}

func (s *usersStore) SaveUser(u user.User) (user.User, error) {
	s.setLoading(true)
	defer s.setLoading(false)

	resp, errResp, err := newRequest[user.User](*s.baseStore).
		Path("users/save").
		Method(http.MethodPost).
		Data(utils.ToJSON(u)).
		Do()

	if err != nil {
		return u, fmt.Errorf("error save user, status: %s, error: %s", err.Error(), errResp.Error)
	}

	result := resp.Data

	s.RequestUsers()
	return result, nil
}

func (s *usersStore) RequestUsers() error {
	s.setLoading(true)
	defer s.setLoading(false)

	resp, errResp, err := newRequest[[]user.User](*s.baseStore).
		Path("users/list").
		Do()

	if err != nil {
		return fmt.Errorf("error request users, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.users = resp.Data

	return nil
}

func (s *usersStore) CheckUniqueLogin(login string, id int64) (bool, error) {
	s.setLoading(true)
	defer s.setLoading(false)

	resp, errResp, err := newRequest[bool](*s.baseStore).
		Path("users/is_exists").
		Param("login", login).
		Param("id", id).
		Do()

	if err != nil {
		return false, fmt.Errorf("error check login unique, status: %s, error: %s", err.Error(), errResp.Error)
	}

	return resp.Data, nil
}

func (s *usersStore) RemoveUser(u user.User) error {
	_, errResp, err := newRequest[user.User](*s.baseStore).
		Path(fmt.Sprintf("users/delete/%d", u.ID)).
		Method(http.MethodDelete).
		Do()

	if err != nil {
		return fmt.Errorf("error save task, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.RequestUsers()

	return nil
}

func (s *usersStore) GetUsers() []user.User {
	return s.users
}

func (s *usersStore) setLoading(v bool) {
	s.loading = v
}

func (s *usersStore) GetUsersLoading() bool {
	return s.loading
}
