package queries

import (
	"foodtastechess/game"
)

type Query interface {
	hasResult() bool
	getResult(dependentQueries ...Query)
	getDependentQueries() []Query
	isExpired(now interface{}) bool
	getExpiration(now interface{}) interface{}
}

type turnNumberQuery struct {
	gameId game.Id

	result game.TurnNumber
}

type moveAtTurnQuery struct {
	gameId game.Id

	result game.AlgebraicMove
}

type boardStateAtTurnQuery struct {
	gameId game.Id

	result game.FEN
}

type validMovesAtTurnQuery struct {
	gameId game.Id

	result game.ValidMoves
}

type unmovedPositionsAtTurnQuery struct {
	gameId game.Id

	result []game.Position
}

type drawOfferStateQuery struct {
	gameId game.Id

	result struct {
		hasOffer bool
		offerer  game.Color
	}
}

type userGamesQuery struct {
	userId string
	result []game.Id
}

/*
type Query struct {
	GameId     game.Id
	QueryType  string
	TurnNumber game.TurnNumber
}
*/

func TurnNumberQuery(id game.Id) Query {
	return turnNumberQuery{
		gameId: id,
	}
}

func BoardAtTurnQuery(id game.Id, turnNumber game.TurnNumber) Query {
	return boardStateAtTurnQuery{
		GameId:     id,
		TurnNumber: turnNumber,
	}
}
