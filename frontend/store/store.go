package store

import (
	v1 "accounter/backend/server/handlers/v1"
	"accounter/config"
	"accounter/domain/user"
	"accounter/tools"
	"fmt"
)

type baseStore struct {
	api  string
	user user.CurrentUser
}

func (b *baseStore) SetUser(user user.CurrentUser) {
	b.user = user
}

func newRequest[R any](s baseStore) tools.Request[v1.Response[R]] {
	params := tools.NewRequest[v1.Response[R]](s.api)

	if s.user.IsAuthorized {
		params = params.Headers(map[string]string{
			"Authorization": s.user.Tokens.AccessToken,
		})
	}

	return params
}

type Store struct {
	baseStore
	mainStore
	tasksStore
	usersStore
}

func NewStore(cfg config.Config) *Store {
	base := baseStore{
		api: fmt.Sprintf("http://localhost:%d/api/v1", cfg.HTTP.Port),
	}

	s := &Store{
		baseStore:  base,
		mainStore:  mainStore{baseStore: base},
		usersStore: usersStore{baseStore: base},
		tasksStore: tasksStore{baseStore: base},
	}

	return s
}
