package store

import (
	"accounter/domain/user"
	"accounter/tools"
)

type mainStore struct {
	baseStore
}

func (s *mainStore) LoginByCredentials(login, password string) error {
	data := tools.Data{
		"login":    login,
		"password": password,
	}

	request := newRequest[user.CurrentUser](s.baseStore).
		Path("login").
		Method("POST").
		Data(data.ToJSON())

	resp, _, err := request.Do()

	if err == nil {
		s.SetUser(resp.Data)
	}

	return err
}

func (s *mainStore) LoginByToken(token string) error {
	headers := map[string]string{
		"Authorization": token,
	}

	request := newRequest[user.CurrentUser](s.baseStore).
		Path("login").
		Headers(headers)

	resp, _, err := request.Do()

	if err == nil {
		s.SetUser(resp.Data)
	}

	return err
}
