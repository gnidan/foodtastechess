package users

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"foodtastechess/config"
)

type Users interface {
	Get(uuid Id) (User, bool)
	GetByAuthId(authId string) (User, bool)
	Save(user *User) error
}

var tablePrefix string = ""

type UsersService struct {
	Config config.DatabaseConfig `inject:"databaseConfig"`

	db gorm.DB
}

func NewUsers() Users {
	return new(UsersService)
}

func (s *UsersService) PostPopulate() error {
	// hook for test-suite, make a global table prefix if our config
	// defines it
	tablePrefix = s.Config.Prefix

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		s.Config.Username, s.Config.Password,
		s.Config.HostAddr, s.Config.Port,
		s.Config.Database,
	)

	db, err := gorm.Open("mysql", dsn)

	db.LogMode(true)
	db.AutoMigrate(&User{})

	s.db = db

	return err
}

func (s *UsersService) Get(uuid Id) (User, bool) {
	user := User{}
	s.db.Where(&User{Uuid: uuid}).First(&user)
	found := (user.Uuid == uuid)
	return user, found
}

func (s *UsersService) GetByAuthId(authId string) (User, bool) {
	user := User{}
	s.db.Where(&User{AuthIdentifier: authId}).First(&user)
	found := (user.AuthIdentifier == authId)
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

// GetAll. Not exposed as part of Users interface
func (s *UsersService) GetAll() []User {
	var users []User
	s.db.Find(&users)
	return users
}

func (s *UsersService) ResetTestDB() {
	if tablePrefix != "test_" {
		log.Error(
			"Cannot reset a database not configured with ConfigTestProvider",
		)
		return
	}
	s.db.DropTable(&User{})
	s.db.AutoMigrate(&User{})
}
