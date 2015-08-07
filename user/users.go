package user

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"foodtastechess/config"
	"foodtastechess/directory"
)

type Users interface {
	Get(uuid string) (User, bool)
	GetByAuthId(authId string) (User, bool)
	Save(user *User) error
}

type UsersService struct {
	Config config.DatabaseConfig `inject:"databaseConfig"`

	db gorm.DB
}

func NewUsers() Users {
	return new(UsersService)
}

func (s *UsersService) PreProvide(provide directory.Provider) error {
	err := provide("databaseConfig",
		config.NewMariaDockerComposeConfig(),
	)

	return err
}

func (s *UsersService) PostPopulate() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		s.Config.Username, s.Config.Password,
		s.Config.HostAddr, s.Config.Port,
		s.Config.Database,
	)

	db, err := gorm.Open("mysql", dsn)

	db.AutoMigrate(&User{})

	s.db = db

	return err
}

func (s *UsersService) Get(uuid string) (User, bool) {
	user := User{}
	s.db.Where(&User{Uuid: uuid}).First(&user)
	found := (user.ID != 0)
	return user, found
}

func (s *UsersService) GetByAuthId(authId string) (User, bool) {
	user := User{}
	s.db.Where(&User{AuthIdentifier: authId}).First(&user)
	found := (user.ID != 0)
	return user, found
}

func (s *UsersService) Save(user *User) error {
	if s.db.NewRecord(*user) {
		s.db.Create(user)
	} else {
		s.db.Save(user)
	}

	return nil
}
