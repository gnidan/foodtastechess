package queries

import (
	"foodtastechess/game"
	"foodtastechess/logger"
)

var log = logger.Log("queries")

type Query interface {
	hasResult() bool
	computeResult(map[Query]Query)
	getDependentQueries() []Query
	isExpired(now interface{}) bool
	getExpiration(now interface{}) interface{}
}

type validMovesAtTurnQuery struct {
	gameId game.Id

	result game.ValidMoves
}

type unmovedPositionsAtTurnQuery struct {
	gameId game.Id

	result []game.Position
}

func TurnNumberQuery(id game.Id) Query {
	return &turnNumberQuery{
		gameId: id,
	}
}

func BoardAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return &boardStateAtTurnQuery{
		gameId:     id,
		turnNumber: turnNumber,
	}
}
