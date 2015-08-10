package events

import (
	"foodtastechess/game"
)

type Events interface {
	Receive(event Event) error

	EventsForGame(gameId game.Id) []Event
	EventsOfTypeForGame(gameId game.Id, eventType string) []Event
	EventsOfTypeForPlayer(userId string, eventType string) []Event
}

type Event struct {
	GameId game.Id
}
