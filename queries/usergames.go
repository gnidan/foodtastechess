package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type userGamesQuery struct {
	PlayerId string

	Answered bool
	Result   []game.Id
}

func (q *userGamesQuery) hasResult() bool {
	return q.Answered
}

func (q *userGamesQuery) getResult() interface{} {
	return q.Result
}

func (q *userGamesQuery) computeResult(queries SystemQueries) {
	activeGames := make(map[game.Id]events.Event)

	gameStarts := queries.getEvents().
		EventsOfTypeForPlayer(q.PlayerId, events.GameStartType)

	gameEnds := queries.getEvents().
		EventsOfTypeForPlayer(q.PlayerId, events.GameEndType)

	for _, event := range gameStarts {
		activeGames[event.GameId] = event
	}

	for _, event := range gameEnds {
		delete(activeGames, event.GameId)
	}

	activeGameIds := []game.Id{}

	for id, _ := range activeGames {
		activeGameIds = append(activeGameIds, id)
	}

	q.Result = activeGameIds
	q.Answered = true
}

func (q *userGamesQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *userGamesQuery) hash() string {
	return fmt.Sprintf("usergames:%v", q.PlayerId)
}
