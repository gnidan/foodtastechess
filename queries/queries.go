package queries

import (
	"foodtastechess/game"
)

type Query struct {
	GameId     game.Id
	QueryType  string
	TurnNumber game.TurnNumber
}

func TurnNumberQuery(id game.Id) Query {
	return Query{
		GameId:    id,
		QueryType: "TurnNumber",
	}
}

func BoardAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return Query{
		GameId:     id,
		QueryType:  "BoardAtTurn",
		TurnNumber: turnNumber,
	}
}

type Answer interface{}
