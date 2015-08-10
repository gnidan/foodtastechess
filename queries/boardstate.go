package queries

import (
	"foodtastechess/game"
)

type boardStateAtTurnQuery struct {
	gameId     game.Id
	turnNumber game.TurnNumber

	result game.FEN
}

func (q *boardStateAtTurnQuery) hasResult() bool {
	return q.result != ""
}

func (q *boardStateAtTurnQuery) computeResult(map[Query]Query) {
}

func (q *boardStateAtTurnQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *boardStateAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *boardStateAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}
