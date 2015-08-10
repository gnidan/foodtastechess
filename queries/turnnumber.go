package queries

import (
	"fmt"

	"foodtastechess/game"
)

type turnNumberQuery struct {
	gameId game.Id

	result game.TurnNumber
}

func (q *turnNumberQuery) hash() string {
	return fmt.Sprintf("turnnumber:%v", q.gameId)
}

func (q *turnNumberQuery) hasResult() bool {
	return q.result != -1
}

func (q *turnNumberQuery) computeResult(queries SystemQueries) {
	moves := queries.getEvents().EventsOfTypeForGame(q.gameId, "move")
	q.result = game.TurnNumber(len(moves))
}

func (q *turnNumberQuery) getDependentQueries() []Query {
	return []Query{}
}
