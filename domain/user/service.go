package user

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

func (us *UserService) GetList() ([]User, error) {
	users, err := us.repo.GetList()

	return users, err
}

func (us *UserService) SaveUser(user *User) error {
	return us.repo.Save(user)
}
