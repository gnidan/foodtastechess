package events

import (
	"foodtastechess/config"
	"foodtastechess/game"
	"foodtastechess/users"
)

type Fixtures struct {
	Config config.FixturesConfig `inject:"fixturesConfig"`
	Events Events                `inject:"events"`
}

func NewFixtures() *Fixtures {
	return new(Fixtures)
}

func (f *Fixtures) Start() error {
	if !f.Config.Enabled {
		return nil
	}

	eventsService := f.Events.(*EventsService)
	eventsService.db.Delete(Event{})

	fixtures := f.gameFixtures()
	for _, event := range fixtures {
		f.Events.Receive(event)
	}

	return nil
}

func (f *Fixtures) Stop() error {
	return nil
}

func (f *Fixtures) gameFixtures() []Event {
	var (
		gameId  game.Id  = 1
		whiteId users.Id = users.Id(f.Config.WhiteId)
		blackId users.Id = users.Id(f.Config.BlackId)
	)

	return []Event{
		NewGameCreateEvent(gameId, whiteId, ""),
		NewGameStartEvent(gameId, whiteId, blackId),
		NewMoveEvent(gameId, 1, "Nb1-c3"),
		NewMoveEvent(gameId, 2, "Nb8-c6"),
		NewMoveEvent(gameId, 3, "Pa2-a4"),
	}
}
