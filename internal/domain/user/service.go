package user

import (
	"accounter/internal/domain/common"
	"accounter/pkg/utils"
	"context"
	"errors"
)

// User service
type UserService struct {
	repo UserRepository
}

// NewUserService creates new UserService
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUsersList get list of User
func (s *UserService) GetUsersList(ctx context.Context, params common.RequestParams) ([]User, error) {
	return s.repo.GetList(ctx, params)
}

// GetUser get single User by specified id
func (s *UserService) GetUser(ctx context.Context, id int64) (User, error) {
	return s.repo.GetOne(ctx, id)
}

// SaveUser create/update User
func (s *UserService) SaveUser(ctx common.Context, user *User) error {
	if utils.IsEmpty(user.ID) {
		return s.repo.Create(ctx, user)
	}

	return s.repo.Update(ctx, user)
}

// DeleteUser delete User by specified id
func (s *UserService) DeleteUser(ctx common.Context, id int64) error {
	user, err := s.GetUser(ctx, id)

	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, user.ID)
}

// CheckUniqueLogin check for User existance by login
func (s *UserService) CheckUniqueLogin(ctx context.Context, login string) (exists bool, err error) {
	_, err = s.repo.GetByCredentials(ctx, login, "")

	if errors.Is(err, ErrUserNotFound) {
		return false, nil

	} else if err == nil {
		exists = true
	}

	return
}
