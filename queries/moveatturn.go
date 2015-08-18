package queries

import (
	"fmt"

	"foodtastechess/game"
)

type moveAtTurnQuery struct {
	GameId     game.Id
	TurnNumber game.TurnNumber

	Answered bool
	Result   game.AlgebraicMove

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *moveAtTurnQuery) hash() string {
	return fmt.Sprintf("move:%v:%v", q.GameId, q.TurnNumber)
}

func (q *moveAtTurnQuery) hasResult() bool {
	return q.Answered
}

func (q *moveAtTurnQuery) getResult() interface{} {
	return q.Result
}

func (q *moveAtTurnQuery) computeResult(queries SystemQueries) {
	moveEvent := queries.getEvents().MoveEventForGameAtTurn(q.GameId, q.TurnNumber)

	q.Result = moveEvent.Move
	q.Answered = true
}

func (q *moveAtTurnQuery) getDependentQueries() []Query {
	return []Query{}
}
