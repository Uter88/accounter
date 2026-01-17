package user

import "context"

// CurrentUser model
type CurrentUser struct {
	context.Context
	User
	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"tokens"`

	IsAuthorized bool `json:"is_authorized"`

	Params struct {
		DateStart int64 `json:"date_start"`
		DateEnd   int64 `json:"date_end"`
	} `json:"params"`
}

// SetToken set authorization tokens
func (u *CurrentUser) SetToken(access, refresh string) {
	u.Tokens.AccessToken = access
	u.Tokens.RefreshToken = refresh
	u.SetAuthorized(true)
}

// SetAuthorized set User authorized flag
func (u *CurrentUser) SetAuthorized(v bool) {
	u.IsAuthorized = v
}

func (u *CurrentUser) WithContext(ctx context.Context) CurrentUser {
	u.Context = ctx
	return *u
}
