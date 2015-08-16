package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type userGamesQuery struct {
	playerId string

	Answered bool
	result   []game.Id
}

func (q *userGamesQuery) hasResult() bool {
	return q.Answered
}

func (q *userGamesQuery) getResult() interface{} {
	return q.result
}

func (q *userGamesQuery) computeResult(queries SystemQueries) {
	activeGames := make(map[game.Id]events.Event)

	gameStarts := queries.getEvents().
		EventsOfTypeForPlayer(q.playerId, events.GameStartType)

	gameEnds := queries.getEvents().
		EventsOfTypeForPlayer(q.playerId, events.GameEndType)

	for _, event := range gameStarts {
		activeGames[event.GameId()] = event
	}

	for _, event := range gameEnds {
		delete(activeGames, event.GameId())
	}

	activeGameIds := []game.Id{}

	for id, _ := range activeGames {
		activeGameIds = append(activeGameIds, id)
	}

	q.result = activeGameIds
	q.Answered = true
}

func (q *userGamesQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *userGamesQuery) hash() string {
	return fmt.Sprintf("usergames:%v", q.playerId)
}
