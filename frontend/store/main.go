package store

import (
	"accounter/domain/user"
	"accounter/frontend/models"
	"accounter/tools"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type mainStore struct {
	*baseStore
}

func (b *mainStore) CheckAuth(ctx app.Context) bool {
	if b.IsAuthorized() {
		return true
	}

	var token string

	if err := ctx.LocalStorage().Get("token", &token); err != nil {
		if err = ctx.LocalStorage().Set("token", &token); err != nil {
			return false
		}
	}

	if token == "" {
		return false
	}

	if err := b.LoginByToken(ctx, token); err != nil {
		return false
	}

	return true
}

func (s *mainStore) LoginByCredentials(ctx app.Context, form models.LoginForm) error {
	data := tools.Data{
		"login":    form.Login,
		"password": form.Password,
	}

	request := newRequest[user.CurrentUser](*s.baseStore).
		Path("login").
		Method("POST").
		Data(data.ToJSON())

	resp, _, err := request.Do()

	if err == nil {
		s.SetUser(ctx, resp.Data, form.IsRemember)
	}

	return err
}

func (s *mainStore) LoginByToken(ctx app.Context, token string) error {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	request := newRequest[user.CurrentUser](*s.baseStore).
		Path("login").
		Headers(headers)

	resp, _, err := request.Do()

	if err == nil {
		s.SetUser(ctx, resp.Data, false)
	}

	return err
}

func (s *mainStore) Logout(ctx app.Context) {
	ctx.SessionStorage().Clear()
	s.user = user.CurrentUser{}
	ctx.Navigate("/login")
}
