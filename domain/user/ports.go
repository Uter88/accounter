package user

type UserRepository interface {
	GetList() ([]User, error)
	GetOne(id int64) (User, error)
	Save(*User) error
	Delete(id int64) error
}
