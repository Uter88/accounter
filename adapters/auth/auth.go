package auth

import (
	"accounter/config"
	"accounter/internal/domain/user"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CurrentUser model
type CurrentUser struct {
	user.User
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

// Authorization service
type AuthService struct {
	repo user.UserRepository
}

// Creates new AuthService
func NewAuthService(repo user.UserRepository) AuthService {
	return AuthService{
		repo: repo,
	}
}

// Authorize CurrentUser by login and password
func (s AuthService) LoginByCredentials(ctx context.Context, login, password string, cfg config.Config) (result CurrentUser, err error) {
	if u, err := s.repo.GetByCredentials(ctx, login, password); err != nil {
		return result, err

	} else {
		result.User = u
		err = createTokens(&result, cfg)

		return result, err
	}
}

// Authorize CurrentUser by JWT token
func (s AuthService) LoginByToken(ctx context.Context, token string, cfg config.Config) (result CurrentUser, err error) {
	if items := strings.Split(token, " "); len(items) == 2 {
		token = items[1]
	}

	if id, err := parseToken(token, cfg.SecretKey); err != nil {
		return result, err

	} else if u, err := s.repo.GetOne(ctx, id); err != nil {
		return result, err

	} else {
		result.User = u
		err = createTokens(&result, cfg)

		return result, err
	}
}

// createTokens creates new JWT token
func createTokens(u *CurrentUser, cfg config.Config) error {
	payload := JWTPayload{UserID: u.ID}

	token, err := payload.GenerateToken(cfg.TokenExpire, cfg.SecretKey)

	if err != nil {
		return err
	}

	u.SetToken(token, "")

	return nil
}

// JWT payload
type JWTPayload struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// JWT tokens pair
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Valid is token expire validation
func (p JWTPayload) Valid() error {
	if time.Now().After(time.Unix(p.ExpiresAt, 0)) {
		return errors.New("token is expired")
	}

	return nil
}

// GenerateToken creates new JWT token
func (p *JWTPayload) GenerateToken(expire time.Duration, secretKey string) (string, error) {
	p.ExpiresAt = time.Now().Add(expire).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, p).SignedString([]byte(secretKey))
}

// Parse token by secret key (salt)
func parseToken(token, secretKey string) (int64, error) {
	claims := JWTPayload{}

	if err := claims.decode(token, secretKey); err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// Decode token
func (p *JWTPayload) decode(token, secretKey string) error {
	_, err := jwt.ParseWithClaims(token, p, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	return err
}
