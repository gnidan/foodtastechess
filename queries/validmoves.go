package queries

import (
	"fmt"

	"foodtastechess/game"
)

type validMovesAtTurnQuery struct {
	GameId     game.Id
	TurnNumber game.TurnNumber

	Answered bool
	Result   []game.MoveRecord

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
	dependentQueries := queries.getDependentQueryLookup(q)
	state := dependentQueries.Lookup(BoardAtTurnQuery(q.GameId, q.TurnNumber)).(*boardStateAtTurnQuery).Result

	q.Result = queries.getGameCalculator().ValidMoves(state)
	q.Answered = true
}

func (q *validMovesAtTurnQuery) getDependentQueries() []Query {
	return []Query{
		BoardAtTurnQuery(q.GameId, q.TurnNumber),
	}
}
