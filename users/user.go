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
	Id                int
	Uuid              Id `sql:"unique_index"`
	Name              string
	AvatarUrl         string
	AuthIdentifier    string `sql:"unique_index"`
	AccessToken       string
	AccessTokenSecret string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u User) TableName() string {
	return fmt.Sprintf("%susers", tablePrefix)
}

func NewId() Id {
	return Id(uuid.NewV4().String())
}
