package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type turnNumberQuery struct {
	GameId game.Id

	Result game.TurnNumber

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *turnNumberQuery) hash() string {
	return fmt.Sprintf("turnnumber:%v", q.GameId)
}

func (q *turnNumberQuery) hasResult() bool {
	return q.Result != -1
}

func (q *turnNumberQuery) getResult() interface{} {
	return q.Result
}

func (q *turnNumberQuery) computeResult(queries SystemQueries) {
	moves := queries.getEvents().EventsOfTypeForGame(q.GameId, events.MoveType)
	q.Result = game.TurnNumber(len(moves))
}

func (q *turnNumberQuery) getDependentQueries() []Query {
	return []Query{}
}
