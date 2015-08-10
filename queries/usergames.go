package queries

import (
	"foodtastechess/game"
)

type userGamesQuery struct {
	gameId game.Id

	result []game.Id
}

func (q *userGamesQuery) hasResult() bool {
	return true
}

func (q *userGamesQuery) computeResult(map[Query]Query) {
}

func (q *userGamesQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *userGamesQuery) isExpired(now interface{}) bool {
	return false
}

func (q *userGamesQuery) getExpiration(now interface{}) interface{} {
	return nil
}
