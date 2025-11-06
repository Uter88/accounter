package user

import "accounter/tools"

// User service
type UserService struct {
	repo UserRepository
}

// NewUserService creates new UserService
func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

// GetUsersList get list of User
func (s *UserService) GetUsersList() ([]User, error) {
	users, err := s.repo.GetList()

	return users, err
}

// SaveUser create/update User
func (s *UserService) SaveUser(user *User) error {
	return s.repo.Save(user)
}

// DeleteUser delete User by id
func (s *UserService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}

// CheckUniqueLogin check for User existance by login
func (s *UserService) CheckUniqueLogin(login string) (exists bool, err error) {
	_, err = s.repo.GetByCredentials(login, "")

	if tools.IsNotFoundError(err) {
		return false, nil
	} else if err == nil {
		exists = true
	}

	return
}
