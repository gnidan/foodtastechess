package user

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"foodtastechess/config"
	"foodtastechess/directory"
)

type Users interface {
	Get(id Id) (User, bool)
	GetByAuthId(authId string) (User, bool)
	Save(user *User) error
}

type UsersService struct {
	Config config.DatabaseConfig `inject:"databaseConfig"`
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
		"%s:%s@tcp(%s:%s)/%s",
		s.Config.Username, s.Config.Password,
		s.Config.HostAddr, s.Config.Port,
		s.Config.Database,
	)

	db, err := gorm.Open("mysql", dsn)
	log.Debug("Got db: %v", db)

	db.AutoMigrate(&User{})

	return err
}

func (s *UsersService) Get(id Id) (User, bool) {
	return User{}, false
}

func (s *UsersService) Save(user User) error {
	return nil
}
