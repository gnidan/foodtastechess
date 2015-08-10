package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type userGamesQuery struct {
	playerId string

	answered bool
	result   []game.Id
}

func (q *userGamesQuery) hasResult() bool {
	return q.answered
}

func (q *userGamesQuery) computeResult(queries SystemQueries) {
	activeGames := make(map[game.Id]events.Event)

	gameStarts := queries.getEvents().EventsOfTypeForPlayer(q.playerId, "start_game")
	gameEnds := queries.getEvents().EventsOfTypeForPlayer(q.playerId, "end_game")

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

	q.result = activeGameIds
	q.answered = true
}

func (q *userGamesQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *userGamesQuery) hash() string {
	return fmt.Sprintf("usergames:%v", q.playerId)
}
