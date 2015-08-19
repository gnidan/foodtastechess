package queries

import (
	"fmt"

	"foodtastechess/game"
)

type validMovesAtTurnQuery struct {
	GameId     game.Id
	TurnNumber game.TurnNumber

	Answered bool
	Result   []game.AlgebraicMove

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *validMovesAtTurnQuery) hash() string {
	return fmt.Sprintf("boardstate:%v:%v", q.GameId, q.TurnNumber)
}

func (q *validMovesAtTurnQuery) hasResult() bool {
	return q.Answered
}

func (q *validMovesAtTurnQuery) getResult() interface{} {
	return q.Result
}

func (q *validMovesAtTurnQuery) computeResult(queries SystemQueries) {
}

func (q *validMovesAtTurnQuery) getDependentQueries() []Query {
	return []Query{
		BoardAtTurnQuery(q.GameId, q.TurnNumber),
	}
}
