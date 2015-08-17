package events

import (
	"foodtastechess/game"
)

type Events interface {
	Receive(event Event) error

	EventsForGame(gameId game.Id) []Event
	EventsOfTypeForGame(gameId game.Id, eventType string) []Event
	EventsOfTypeForPlayer(userId string, eventType string) []Event
	MoveEventForGameAtTurn(gameId game.Id, turnNumber game.TurnNumber) MoveEvent
}

type EventSubscriber interface {
	Receive(event Event) error
}

type Event interface {
	GameId() game.Id
}

type MoveEvent struct {
	gameId     game.Id
	TurnNumber game.TurnNumber
	Move       game.AlgebraicMove
}

func (e MoveEvent) GameId() game.Id {
	return e.gameId
}

type GameStartEvent struct {
	gameId game.Id
}

func (e GameStartEvent) GameId() game.Id {
	return e.gameId
}

type GameEndEvent struct {
	gameId game.Id
}

func (e GameEndEvent) GameId() game.Id {
	return e.gameId
}

type DrawOfferEvent struct {
	gameId game.Id
	Color  game.Color
}

func (e DrawOfferEvent) GameId() game.Id {
	return e.gameId
}

type DrawOfferResponseEvent struct {
	gameId game.Id
	accept bool
}

func (e DrawOfferResponseEvent) GameId() game.Id {
	return e.gameId
}

const (
	MoveType              = "move"
	GameStartType         = "game:start"
	GameEndType           = "game:end"
	DrawOfferType         = "offer:create"
	DrawOfferResponseType = "offer:respond"
)

func NewMoveEvent(gameId game.Id, turnNumber game.TurnNumber, move game.AlgebraicMove) Event {
	return MoveEvent{gameId, turnNumber, move}
}

func NewGameStartEvent(gameId game.Id) Event {
	return GameStartEvent{gameId}
}

func NewGameEndEvent(gameId game.Id) Event {
	return GameEndEvent{gameId}
}

func NewDrawOfferEvent(gameId game.Id, color game.Color) Event {
	return DrawOfferEvent{gameId, color}
}

func NewDrawOfferResponseEvent(gameId game.Id, accept bool) Event {
	return DrawOfferResponseEvent{gameId, accept}
}
