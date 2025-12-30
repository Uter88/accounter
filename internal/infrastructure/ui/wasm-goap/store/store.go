package store

import (
	"accounter/config"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"accounter/pkg/tools"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type baseStore struct {
	api  string
	user user.CurrentUser
}

func (b *baseStore) SetUser(ctx app.Context, user user.CurrentUser, remember bool) {
	b.user = user

	b.storeAccessToken(ctx, user.Tokens.AccessToken)

	if remember {
		b.storeAuthData(ctx)
	}
}

func (b *baseStore) storeAccessToken(ctx app.Context, token string) {
	ctx.SessionStorage().Set("access_token", token)
}

func (b *baseStore) loadAccessToken(ctx app.Context) (token string, ok bool) {
	if err := ctx.SessionStorage().Get("access_token", &token); err != nil {
		return
	}

	ok = token != ""

	return
}

func (b *baseStore) storeAuthData(ctx app.Context) {
	authData := fmt.Sprintf("%s:%s", b.user.Login, b.user.Password)
	authData = base64.StdEncoding.EncodeToString([]byte(authData))
	ctx.LocalStorage().Set("auth_data", authData)
}

func (b *baseStore) LoadAuthData(ctx app.Context) (login, password string, ok bool) {
	var authData string

	if err := ctx.LocalStorage().Get("auth_data", &authData); err != nil {
		return
	} else if res, err := base64.StdEncoding.DecodeString(authData); err != nil {
		return
	} else {
		authData = string(res)
	}

	if items := strings.Split(authData, ":"); len(items) != 2 {
		return
	} else {
		login = items[0]
		password = items[1]
		ok = true
	}

	return
}

func (b *baseStore) GetUser() user.CurrentUser {
	return b.user
}

func (b *baseStore) IsAuthorized() bool {
	return b.user.IsAuthorized
}

func newRequest[R any](s baseStore) tools.Request[tools.Response[R]] {
	params := tools.NewRequest[tools.Response[R]](s.api)

	if s.user.IsAuthorized {
		params = params.Headers(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", s.user.Tokens.AccessToken),
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
		tasksStore: tasksStore{
			baseStore: base,
			params:    task.NewTaskParams(),
		},
	}

	return s
}
