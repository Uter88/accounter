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
func (s *AuthService) LoginByCredentials(login, password string, cfg config.Config) (result user.CurrentUser, err error) {

	if u, err := s.repo.GetByCredentials(login, password); err != nil {
		return result, err
	} else {
		result.User = u
		payload := JWTPayload{UserID: u.ID}
		token, err := payload.GenerateToken(time.Hour*24*31, cfg.SecretKey)

		if err != nil {
			return result, err
		}

		result.SetToken(token, "")
	}

	return
}

// Authorize CurrentUser by JWT token
func (s *AuthService) LoginByToken(c *gin.Context, cfg config.Config) (result user.CurrentUser, err error) {
	token := c.GetHeader("Authorization")

	if id, err := parseToken(token, cfg.SecretKey); err != nil {
		return result, err

	} else if u, err := s.repo.GetOne(id); err != nil {
		return result, err

	} else {
		result.User = u
		result.SetToken(token, "")

		return result, err
	}
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

// Decode token
func (p *JWTPayload) decode(token, secretKey string) error {
	_, err := jwt.ParseWithClaims(token, p, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	return err
}
