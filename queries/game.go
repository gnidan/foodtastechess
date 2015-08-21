package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type gameQuery struct {
	GameId game.Id

	Answered bool
	Result   bool

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *gameQuery) hasResult() bool {
	return q.Answered
}

func (q *gameQuery) getResult() interface{} {
	return q.Result
}

func (q *gameQuery) computeResult(queries SystemQueries) {
	gameStarts := queries.getEvents().
		EventsOfTypeForGame(q.GameId, events.GameCreateType)

	q.Answered = true
	q.Result = len(gameStarts) > 0
}

func (q *gameQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *gameQuery) hash() string {
	return fmt.Sprintf("game:%v", q.GameId)
}
