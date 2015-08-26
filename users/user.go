package users

import (
	"database/sql/driver"
	"fmt"
	"github.com/satori/go.uuid"
	"time"

	"foodtastechess/logger"
)

var log = logger.Log("user")

type Id string

func (u *Id) Scan(value interface{}) error {
	*u = Id(value.([]byte))
	return nil
}

func (u Id) Value() (driver.Value, error) {
	return string(u), nil
}

type User struct {
	Id                int `json:"-"`
	Uuid              Id  `sql:"unique_index" json:"-"`
	Name              string
	AvatarUrl         string
	AuthIdentifier    string `sql:"unique_index" json:"-"`
	AccessToken       string `json:"-"`
	AccessTokenSecret string `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) TableName() string {
	return fmt.Sprintf("%susers", tablePrefix)
}

func NewId() Id {
	return Id(uuid.NewV4().String())
}
