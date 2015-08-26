package commands

import (
	"foodtastechess/game"
	"foodtastechess/users"
)

type context struct {
	name        string
	userId      users.Id
	gameId      game.Id
	move        game.AlgebraicMove
	colorChoice game.Color
	accept      bool
}
