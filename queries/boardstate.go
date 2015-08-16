package queries

import (
	"fmt"

	"foodtastechess/game"
)

type boardStateAtTurnQuery struct {
	GameId     game.Id
	TurnNumber game.TurnNumber

	Result game.FEN

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *boardStateAtTurnQuery) hash() string {
	return fmt.Sprintf("boardstate:%v:%v", q.GameId, q.TurnNumber)
}

func (q *boardStateAtTurnQuery) hasResult() bool {
	return q.Result != ""
}

func (q *boardStateAtTurnQuery) getResult() interface{} {
	return q.Result
}

func (q *boardStateAtTurnQuery) computeResult(queries SystemQueries) {
	if q.TurnNumber == 0 {
		q.Result = queries.getGameCalculator().StartingFEN()
		return
	}

	dependentQueries := queries.getDependentQueryLookup(q)
	log.Debug("%v", dependentQueries)
	lastPosition := dependentQueries.Lookup(BoardAtTurnQuery(q.GameId, q.TurnNumber-1)).(*boardStateAtTurnQuery).Result

	lastMove := dependentQueries.Lookup(MoveAtTurnQuery(q.GameId, q.TurnNumber)).(*moveAtTurnQuery).Result

	q.Result = queries.getGameCalculator().AfterMove(lastPosition, lastMove)
}

func (q *boardStateAtTurnQuery) getDependentQueries() []Query {
	if q.TurnNumber == 0 {
		return []Query{}
	} else {
		return []Query{
			BoardAtTurnQuery(q.GameId, q.TurnNumber-1),
			MoveAtTurnQuery(q.GameId, q.TurnNumber),
		}
	}
}

func (q *boardStateAtTurnQuery) GoString() string {
	return fmt.Sprintf(
		"BoardAtTurn(%d, game=%d)", q.TurnNumber, q.GameId,
	)
}
