package core

import (
	"accounter/config"
	"accounter/domain/user"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CurrentUser struct {
	user.User
}

func NewAuthService(repo user.UserRepository) AuthService {
	return AuthService{
		repo: repo,
	}
}

type AuthService struct {
	repo user.UserRepository
}

func (m *AuthService) TokenAuthorization(c *gin.Context, cfg config.Config) (*CurrentUser, error) {
	token := c.GetHeader("Authorization")

	if id, err := parseToken(token, cfg.SecretKey); err != nil {
		return nil, err

	} else if user, err := m.repo.GetOne(id); err != nil {
		return nil, err

	} else {
		return &CurrentUser{User: user}, nil
	}
}

type JWTPayload struct {
	UserID int64 `json:"user_key"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (p JWTPayload) Valid() error {
	if time.Now().After(time.Unix(p.ExpiresAt, 0)) {
		return errors.New("token is expired")
	}

	return nil
}

func (p *JWTPayload) GenerateToken(expire time.Duration, secretKey string) (string, error) {
	p.ExpiresAt = time.Now().Add(expire).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, p).SignedString([]byte(secretKey))
}

func parseToken(data, secretKey string) (int64, error) {
	token := strings.SplitN(data, " ", 2)

	if len(token) != 2 {
		return 0, errors.New("miss auth token")
	}

	claims := JWTPayload{}

	if err := claims.decode(token[1], secretKey); err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

func (p *JWTPayload) decode(token, secretKey string) error {
	_, err := jwt.ParseWithClaims(token, p, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	return err
}
