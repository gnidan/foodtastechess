package queries

import (
	"foodtastechess/game"
)

type moveAtTurnQuery struct {
	gameId game.Id

	result game.AlgebraicMove
}

func (q *moveAtTurnQuery) hasResult() bool {
	return true
}

func (q *moveAtTurnQuery) computeResult(map[Query]Query) {
}

func (q *moveAtTurnQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *moveAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *moveAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}
