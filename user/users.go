package user

type Users interface {
	Get(id Id) (User, bool)
	Save(user User) error
}

type UsersService struct {
}

func NewUsers() Users {
	return new(UsersService)
}

func (s *UsersService) Get(id Id) (User, bool) {
	return User{}, false
}

func (s *UsersService) Save(user User) error {
	return nil
}
