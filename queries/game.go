package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type gameQuery struct {
	GameId game.Id

	Answered bool
	Result   GameStatus

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
	q.Answered = true

	exists := func(e events.EventType) bool {
		return len(
			queries.getEvents().
				EventsOfTypeForGame(q.GameId, e),
		) > 0
	}

	if !exists(events.GameCreateType) {
		q.Result = GameStatusNull
	} else if !exists(events.GameStartType) {
		q.Result = GameStatusCreated
	} else if !exists(events.GameEndType) {
		q.Result = GameStatusStarted
	} else {
		q.Result = GameStatusEnded
	}
}

func (q *gameQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *gameQuery) hash() string {
	return fmt.Sprintf("game:%v", q.GameId)
}
