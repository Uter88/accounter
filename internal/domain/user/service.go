package user

import (
	"accounter/internal/domain/shared"
	"accounter/pkg/tools"
	"context"
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
func (s *UserService) GetUsersList(ctx context.Context) ([]User, error) {
	return s.repo.GetList(ctx)
}

// GetUser get single User by specified id
func (s *UserService) GetUser(ctx context.Context, id int64) (User, error) {
	return s.repo.GetOne(ctx, id)
}

// SaveUser create/update User
func (s *UserService) SaveUser(ctx shared.Context, user *User) error {
	if tools.IsEmpty(user.ID) {
		return s.repo.Create(ctx, user)
	}

	return s.repo.Update(ctx, user)
}

// DeleteUser delete User
func (s *UserService) DeleteUser(ctx shared.Context, user *User) error {
	return s.repo.Delete(ctx, user.ID)
}

// CheckUniqueLogin check for User existance by login
func (s *UserService) CheckUniqueLogin(ctx context.Context, login string) (exists bool, err error) {
	_, err = s.repo.GetByCredentials(ctx, login, "")

	if tools.IsNotFoundError(err) {
		return false, nil

	} else if err == nil {
		exists = true
	}

	return
}
