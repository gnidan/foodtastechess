package user

import (
	"foodtastechess/logger"
	"time"
)

var log = logger.Log("user")

type Id string

type User struct {
	ID             int
	Uuid           Id
	Name           string
	AvatarUrl      string
	AuthIdentifier string

	CreatedAt time.Time
	UpdatedAt time.Time
}
