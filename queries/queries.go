package queries

import (
	"time"

	"foodtastechess/game"
	"foodtastechess/logger"
)

var log = logger.Log("queries")

type Query interface {
	hasResult() bool
	getResult() interface{}

	getDependentQueries() []Query
	computeResult(sqs SystemQueries)

	isExpired(now interface{}) bool
	getExpiration(now interface{}) interface{}

	hash() string
}

type queryRecord struct {
	Hash       string
	ComputedAt time.Time
}

type validMovesAtTurnQuery struct {
	GameId game.Id

	Result game.ValidMoves
}

func TurnNumberQuery(id game.Id) Query {
	return &turnNumberQuery{
		GameId: id,
		Result: -1,
	}
}

func BoardAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return &boardStateAtTurnQuery{
		GameId:     id,
		TurnNumber: turnNumber,
	}
}

func MoveAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return &moveAtTurnQuery{
		GameId:     id,
		TurnNumber: turnNumber,
	}
}

func UserGamesQuery(playerId string) Query {
	return &userGamesQuery{
		PlayerId: playerId,
	}
}

func DrawOfferStateQuery(gameId game.Id) Query {
	return &drawOfferStateQuery{
		GameId: gameId,
	}
}

func GamePlayersQuery(id game.Id) Query {
	return &gamePlayersQuery{
		GameId: id,
	}
}
