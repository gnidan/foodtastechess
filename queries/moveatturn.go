package queries

import (
	"fmt"

	"foodtastechess/game"
)

type moveAtTurnQuery struct {
	gameId     game.Id
	turnNumber game.TurnNumber

	answered bool
	result   game.AlgebraicMove
}

func (q *moveAtTurnQuery) hash() string {
	return fmt.Sprintf("move:%v:%v", q.gameId, q.turnNumber)
}

func (q *moveAtTurnQuery) hasResult() bool {
	return q.answered
}

func (q *moveAtTurnQuery) getResult() interface{} {
	return q.result
}

func (q *moveAtTurnQuery) computeResult(queries SystemQueries) {
	moveEvent := queries.getEvents().MoveEventForGameAtTurn(q.gameId, q.turnNumber)

	q.result = moveEvent.Move
	q.answered = true
}

func (q *moveAtTurnQuery) getDependentQueries() []Query {
	return []Query{}
}
