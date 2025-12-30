package store

import (
	"accounter/internal/domain/user"
	"accounter/internal/infrastructure/ui/wasm-goap/models"
	"accounter/pkg/tools"
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

	token, ok := b.loadAccessToken(ctx)

	if !ok {
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

	resp, errResp, err := newRequest[user.CurrentUser](*s.baseStore).
		Path("login").
		Method("POST").
		Data(data.ToJSON()).
		Do()

	if err != nil {
		return fmt.Errorf("auth error, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.SetUser(ctx, resp.Data, form.IsRemember)

	return nil
}

func (s *mainStore) LoginByToken(ctx app.Context, token string) error {
	resp, errResp, err := newRequest[user.CurrentUser](*s.baseStore).
		Path("login").
		Header("Authorization", fmt.Sprintf("Bearer %s", token)).
		Do()

	if err != nil {
		return fmt.Errorf("auth error, status: %s, error: %s", err.Error(), errResp.Error)
	}

	s.SetUser(ctx, resp.Data, false)

	return nil
}

func (s *mainStore) Logout(ctx app.Context) {
	ctx.SessionStorage().Clear()
	s.user = user.CurrentUser{}
	ctx.Navigate("/login")
}
