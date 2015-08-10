package queries

import (
	"fmt"

	"foodtastechess/game"
)

type moveAtTurnQuery struct {
	gameId     game.Id
	turnNumber game.TurnNumber

	result game.AlgebraicMove
}

func (q *moveAtTurnQuery) hash() string {
	return fmt.Sprintf("move:%v:%v", q.gameId, q.turnNumber)
}

func (q *moveAtTurnQuery) hasResult() bool {
	return true
}

func (q *moveAtTurnQuery) computeResult(queries SystemQueries) {
}

func (q *moveAtTurnQuery) getDependentQueries() []Query {
	return []Query{}
}
