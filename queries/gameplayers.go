package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/users"
)

type gamePlayersQuery struct {
	GameId game.Id

	Answered bool
	Result   map[game.Color]users.Id
}

func (q *gamePlayersQuery) hasResult() bool {
	return q.Answered
}

func (q *gamePlayersQuery) getResult() interface{} {
	return q.Result
}

func (q *gamePlayersQuery) computeResult(queries SystemQueries) {
	gameStart := queries.getEvents().
		EventsOfTypeForGame(q.GameId, events.GameStartType)[0]

	q.Result = make(map[game.Color]users.Id)
	q.Result[game.White] = gameStart.WhiteId
	q.Result[game.Black] = gameStart.BlackId

	q.Answered = true
}

func (q *gamePlayersQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *gamePlayersQuery) hash() string {
	return fmt.Sprintf("gameplayers:%v", q.GameId)
}
