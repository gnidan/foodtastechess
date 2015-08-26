package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type gameEndQuery struct {
	GameId game.Id

	Answered bool
	Result   GameEnd

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *gameEndQuery) hasResult() bool {
	return q.Answered
}

func (q *gameEndQuery) getResult() interface{} {
	return q.Result
}

func (q *gameEndQuery) computeResult(queries SystemQueries) {
	q.Answered = true

	gameEnds := queries.getEvents().
		EventsOfTypeForGame(q.GameId, events.GameEndType)
	if len(gameEnds) == 0 {
		q.Result = GameEnd{Occurred: false}
		return
	}

	gameEnd := gameEnds[0]

	q.Result = GameEnd{
		Occurred: true,
		Reason:   gameEnd.Reason,
		Winner:   gameEnd.Winner,
	}
}

func (q *gameEndQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *gameEndQuery) hash() string {
	return fmt.Sprintf("gameend:%v", q.GameId)
}
