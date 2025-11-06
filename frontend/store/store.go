package store

import (
	v1 "accounter/backend/server/handlers/v1"
	"accounter/config"
	"accounter/domain/user"
	"accounter/tools"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type baseStore struct {
	api  string
	user user.CurrentUser
}

func (b *baseStore) SetUser(ctx app.Context, user user.CurrentUser, remember bool) {
	b.user = user

	if remember {
		ctx.LocalStorage().Set("token", user.Tokens.AccessToken)
	} else {
		ctx.SessionStorage().Set("token", user.Tokens.AccessToken)
	}
}

func (b *baseStore) GetUser() user.CurrentUser {
	return b.user
}

func (b *baseStore) IsAuthorized() bool {
	return b.user.IsAuthorized
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
	*baseStore
	mainStore
	tasksStore
	usersStore
}

func NewStore(cfg config.Config) *Store {
	base := &baseStore{
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
