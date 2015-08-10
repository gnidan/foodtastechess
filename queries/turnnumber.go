package queries

import (
	"foodtastechess/game"
)

type turnNumberQuery struct {
	gameId game.Id

	result game.TurnNumber
}

func (q *turnNumberQuery) hasResult() bool {
	return true
}

func (q *turnNumberQuery) computeResult(map[Query]Query) {
}

func (q *turnNumberQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *turnNumberQuery) isExpired(now interface{}) bool {
	return false
}

func (q *turnNumberQuery) getExpiration(now interface{}) interface{} {
	return nil
}
