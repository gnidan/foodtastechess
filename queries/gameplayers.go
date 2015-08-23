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

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *gamePlayersQuery) hasResult() bool {
	return q.Answered
}

func (q *gamePlayersQuery) getResult() interface{} {
	return q.Result
}

func (q *gamePlayersQuery) computeResult(queries SystemQueries) {
	q.Answered = true

	dependentQueries := queries.getDependentQueryLookup(q)
	status := dependentQueries.Lookup(GameQuery(q.GameId)).getResult().(GameStatus)

	if status == GameStatusNull {
		return
	}

	var whiteId users.Id
	var blackId users.Id

	if status == GameStatusCreated {
		gameCreate := queries.getEvents().
			EventsOfTypeForGame(q.GameId, events.GameCreateType)[0]
		whiteId = gameCreate.WhiteId
		blackId = gameCreate.BlackId
	} else {
		gameStart := queries.getEvents().
			EventsOfTypeForGame(q.GameId, events.GameStartType)[0]
		whiteId = gameStart.WhiteId
		blackId = gameStart.BlackId

	}

	q.Result = make(map[game.Color]users.Id)
	q.Result[game.White] = whiteId
	q.Result[game.Black] = blackId
}

func (q *gamePlayersQuery) getDependentQueries() []Query {
	return []Query{
		GameQuery(q.GameId),
	}
}

func (q *gamePlayersQuery) hash() string {
	return fmt.Sprintf("gameplayers:%v", q.GameId)
}
