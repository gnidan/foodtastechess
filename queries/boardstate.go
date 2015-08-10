package queries

import (
	"fmt"

	"foodtastechess/game"
)

type boardStateAtTurnQuery struct {
	gameId     game.Id
	turnNumber game.TurnNumber

	result game.FEN
}

func (q *boardStateAtTurnQuery) hash() string {
	return fmt.Sprintf("boardstate:%v:%v", q.gameId, q.turnNumber)
}

func (q *boardStateAtTurnQuery) hasResult() bool {
	return q.result != ""
}

func (q *boardStateAtTurnQuery) computeResult(queries SystemQueries) {
	if q.turnNumber == 0 {
		q.result = queries.getGameCalculator().StartingFEN()
		return
	}

	dependentQueries := queries.getDependentQueryLookup(q)
	log.Debug("%v", dependentQueries)
	lastPosition := dependentQueries.Lookup(BoardAtTurnQuery(q.gameId, q.turnNumber-1)).(*boardStateAtTurnQuery).result

	lastMove := dependentQueries.Lookup(MoveAtTurnQuery(q.gameId, q.turnNumber)).(*moveAtTurnQuery).result

	q.result = queries.getGameCalculator().AfterMove(lastPosition, lastMove)
}

func (q *boardStateAtTurnQuery) getDependentQueries() []Query {
	if q.turnNumber == 0 {
		return []Query{}
	} else {
		return []Query{
			BoardAtTurnQuery(q.gameId, q.turnNumber-1),
			MoveAtTurnQuery(q.gameId, q.turnNumber),
		}
	}
}

func (q *boardStateAtTurnQuery) GoString() string {
	return fmt.Sprintf(
		"BoardAtTurn(%d, game=%d)", q.turnNumber, q.gameId,
	)
}
