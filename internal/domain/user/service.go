package user

import (
	"accounter/pkg/tools"
	"context"
)

// User service
type UserService struct {
	repo UserRepository
}

// NewUserService creates new UserService
func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

// GetUsersList get list of User
func (s UserService) GetUsersList(ctx context.Context) ([]User, error) {
	result, err := s.repo.GetList(ctx)

	return result, err
}

// SaveUser create/update User
func (s UserService) SaveUser(ctx context.Context, user *User) error {
	if tools.IsEmpty(user.ID) {
		return s.repo.Insert(ctx, user)
	}

	return s.repo.Update(ctx, user)
}

// DeleteUser delete User by id
func (s UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// CheckUniqueLogin check for User existance by login
func (s UserService) CheckUniqueLogin(ctx context.Context, login string) (exists bool, err error) {
	_, err = s.repo.GetByCredentials(ctx, login, "")

	if tools.IsNotFoundError(err) {
		return false, nil
	} else if err == nil {
		exists = true
	}

	return
}
