package queries

import (
	"foodtastechess/game"
)

type Query struct {
	GameId     game.Id
	QueryType  string
	TurnNumber game.TurnNumber
}

type Answer interface{}
