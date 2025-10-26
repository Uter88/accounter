package store

import (
	"accounter/config"
	"fmt"
)

type baseStore struct {
	api string
}

func (s *baseStore) getApi(path string) string {
	return fmt.Sprintf("%s/%s", s.api, path)
}

type Store struct {
	mainStore
	tasksStore
	usersStore
}

func NewStore(cfg config.Config) *Store {
	base := baseStore{
		api: fmt.Sprintf("http://localhost:%d/api/v1", cfg.HTTP.Port),
	}

	s := &Store{
		mainStore:  mainStore{baseStore: base},
		usersStore: usersStore{baseStore: base},
		tasksStore: tasksStore{baseStore: base},
	}

	return s
}
