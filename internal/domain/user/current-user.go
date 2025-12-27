package user

// CurrentUser model
type CurrentUser struct {
	User
	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"tokens"`

	IsAuthorized bool `json:"is_authorized"`
}

func (u *CurrentUser) SetToken(access, refresh string) {
	u.Tokens.AccessToken = access
	u.Tokens.RefreshToken = refresh
	u.SetAuthorized(true)
}

func (u *CurrentUser) SetAuthorized(v bool) {
	u.IsAuthorized = v
}
