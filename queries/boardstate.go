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

func (q *boardStateAtTurnQuery) hasResult() bool {
	return q.result != ""
}

func (q *boardStateAtTurnQuery) computeResult(queries SystemQueries) {
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

func (q *boardStateAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *boardStateAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}

func (q *boardStateAtTurnQuery) GoString() string {
	return fmt.Sprintf(
		"BoardAtTurn(%d, game=%d)", q.turnNumber, q.gameId,
	)
}
