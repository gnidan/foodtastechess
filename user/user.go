package user

import (
	"fmt"
	"github.com/satori/go.uuid"
	"time"

	"foodtastechess/logger"
)

var log = logger.Log("user")

type User struct {
	ID                int
	Uuid              string `sql:"unique_index"`
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

func NewId() string {
	return uuid.NewV4().String()
}
