package store

import (
	"accounter/config"
	"accounter/internal/domain/common"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"accounter/internal/infrastructure/ui/wasm-goap/models"
	"accounter/pkg/utils"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type baseStore struct {
	api  string
	user user.CurrentUser
	ws   *models.WebsocketClient
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

func newRequest[R any](s baseStore) utils.Request[common.Response[R]] {
	params := utils.NewRequest[common.Response[R]](s.api)

	if s.user.IsAuthorized {
		params = params.Headers(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", s.user.Tokens.AccessToken),
		})
	}

	return params
}

// Store storage for entities
type Store struct {
	*baseStore
	mainStore
	tasksStore
	usersStore
}

// NewStore creates new Store
func NewStore(cfg config.Config, ws *models.WebsocketClient) *Store {
	base := &baseStore{
		api: fmt.Sprintf("http://localhost:%d/api/v1", cfg.HTTP.Port),
		ws:  ws,
	}

	params := task.NewTaskParams()

	s := &Store{
		baseStore:  base,
		mainStore:  mainStore{baseStore: base},
		usersStore: usersStore{baseStore: base},
		tasksStore: tasksStore{
			baseStore: base,
			params:    &params,
		},
	}

	return s
}
