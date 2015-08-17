package events

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"

	"foodtastechess/config"
	"foodtastechess/game"
	"foodtastechess/logger"
)

var tablePrefix string = ""

type Events interface {
	Receive(event Event) error

	EventsForGame(gameId game.Id) []Event
	EventsOfTypeForGame(gameId game.Id, eventType EventType) []Event
	EventsOfTypeForPlayer(userId string, eventType EventType) []Event
	MoveEventForGameAtTurn(gameId game.Id, turnNumber game.TurnNumber) Event
}

type EventSubscriber interface {
	Receive(event Event) error
}

type EventsService struct {
	Config     config.DatabaseConfig `inject:"databaseConfig"`
	Subscriber EventSubscriber       `inject:"eventSubscriber"`

	log *logging.Logger
	db  gorm.DB
}

func (s *EventsService) PostPopulate() error {
	// hook for test-suite, make a global table prefix if our config
	// defines it
	tablePrefix = s.Config.Prefix

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		s.Config.Username, s.Config.Password,
		s.Config.HostAddr, s.Config.Port,
		s.Config.Database,
	)

	db, err := gorm.Open("mysql", dsn)

	db.AutoMigrate(&Event{})

	db.LogMode(true)

	s.db = db

	return err
}

func NewEvents() Events {
	service := new(EventsService)
	service.log = logger.Log("events")
	return service
}

func (s *EventsService) Receive(event Event) error {
	var partial Event = event

	s.db.Create(&event)

	s.Subscriber.Receive(partial)

	return nil
}

func (s *EventsService) EventsForGame(gameId game.Id) []Event {
	var events []Event
	s.db.Where(&Event{GameId: gameId}).Find(&events)
	return events
}

func (s *EventsService) EventsOfTypeForGame(gameId game.Id, eventType EventType) []Event {
	var events []Event
	s.db.Where(&Event{GameId: gameId, Type: eventType}).Find(&events)
	return events
}

func (s *EventsService) EventsOfTypeForPlayer(userId string, eventType EventType) []Event {
	var events []Event
	s.db.
		Where(&Event{Type: eventType, WhiteId: userId}).
		Or(&Event{Type: eventType, BlackId: userId}).
		Find(&events)
	return events
}

func (s *EventsService) MoveEventForGameAtTurn(gameId game.Id, turnNumber game.TurnNumber) Event {
	var event Event
	s.db.
		Where(
		&Event{
			Type:       MoveType,
			GameId:     gameId,
			TurnNumber: turnNumber,
		}).
		First(&event)
	return event
}
