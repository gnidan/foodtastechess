package user

import (
	"foodtastechess/logger"
)

var log = logger.Log("user")

type Id string

type User struct {
	Id        Id
	NickName  string
	AvatarUrl string
}
