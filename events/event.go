package events

import (
	"database/sql/driver"
	"fmt"
	"time"

	"foodtastechess/game"
	"foodtastechess/users"
)

type Event struct {
	Id          int
	Type        EventType       `sql:"index"`
	GameId      game.Id         `sql:"index"`
	TurnNumber  game.TurnNumber `sql:"index"`
	WhiteId     users.Id        `sql:"index"`
	BlackId     users.Id        `sql:"index"`
	Move        game.AlgebraicMove
	Offerer     game.Color
	OfferAccept bool

	CreatedAt time.Time
}

func (e Event) TableName() string {
	return fmt.Sprintf("%sevents", tablePrefix)
}

type EventType string

func (u *EventType) Scan(value interface{}) error {
	*u = EventType(value.([]byte))
	return nil
}

func (u EventType) Value() (driver.Value, error) {
	return string(u), nil
}

const (
	MoveType              EventType = "move"
	GameCreateType        EventType = "game:create"
	GameStartType         EventType = "game:start"
	GameEndType           EventType = "game:end"
	DrawOfferType         EventType = "offer:create"
	DrawOfferResponseType EventType = "offer:respond"
	// don't forget to add to queries/buffer.go if necessary
)

func NewMoveEvent(gameId game.Id, turnNumber game.TurnNumber, move game.AlgebraicMove) Event {
	event := new(Event)
	event.Type = MoveType
	event.GameId = gameId
	event.TurnNumber = turnNumber
	event.Move = move
	return *event
}

func NewGameCreateEvent(gameId game.Id, whiteId, blackId users.Id) Event {
	event := new(Event)
	event.Type = GameCreateType
	event.GameId = gameId
	event.WhiteId = whiteId
	event.BlackId = blackId
	return *event
}

func NewGameStartEvent(gameId game.Id, whiteId, blackId users.Id) Event {
	event := new(Event)
	event.Type = GameStartType
	event.GameId = gameId
	event.WhiteId = whiteId
	event.BlackId = blackId
	return *event
}

func NewGameEndEvent(gameId game.Id, whiteId, blackId users.Id) Event {
	event := new(Event)
	event.Type = GameEndType
	event.GameId = gameId
	event.WhiteId = whiteId
	event.BlackId = blackId
	return *event
}

func NewDrawOfferEvent(gameId game.Id, color game.Color) Event {
	event := new(Event)
	event.Type = DrawOfferType
	event.GameId = gameId
	event.Offerer = color
	return *event
}

func NewDrawOfferResponseEvent(gameId game.Id, accept bool) Event {
	event := new(Event)
	event.Type = DrawOfferResponseType
	event.GameId = gameId
	event.OfferAccept = accept
	return *event
}
