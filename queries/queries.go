package queries

import (
	"foodtastechess/game"
	"foodtastechess/logger"
)

var log = logger.Log("queries")

type Query interface {
	hasResult() bool

	getDependentQueries() []Query
	computeResult(sqs SystemQueries)

	isExpired(now interface{}) bool
	getExpiration(now interface{}) interface{}

	hash() string
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
		result: -1,
	}
}

func BoardAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return &boardStateAtTurnQuery{
		gameId:     id,
		turnNumber: turnNumber,
	}
}

func MoveAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return &moveAtTurnQuery{
		gameId:     id,
		turnNumber: turnNumber,
	}
}

func UserGamesQuery(playerId string) Query {
	return &userGamesQuery{
		playerId: playerId,
	}
}

func DrawOfferStateQuery(gameId game.Id) Query {
	return &drawOfferStateQuery{
		gameId: gameId,
	}
}
