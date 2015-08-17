package queries

import (
	"github.com/op/go-logging"
	"reflect"

	"foodtastechess/events"
	"foodtastechess/logger"
)

type QueryBuffer struct {
	log           *logging.Logger
	queries       chan Query
	SystemQueries SystemQueries `inject:"systemQueries"`
	stopChan      chan bool
}

func NewQueryBuffer() events.EventSubscriber {
	buffer := new(QueryBuffer)
	buffer.log = logger.Log("querybuffer")
	buffer.queries = make(chan Query, 100)
	buffer.stopChan = make(chan bool, 1)
	return buffer
}

func (b *QueryBuffer) Start() error {
	b.log.Notice("Listening for Queries")
	go b.Process()
	return nil
}

func (b *QueryBuffer) Process() {
	for {
		select {
		case query := <-b.queries:
			b.log.Info("Got query")
			b.SystemQueries.computeAnswer(query, false)
		case <-b.stopChan:
			return
		}
	}
	b.log.Info("QueryBuffer stopped")
}

func (b *QueryBuffer) Stop() error {
	b.log.Notice("Stopping QueryBuffer")
	b.stopChan <- true
	return nil
}

func (b *QueryBuffer) Receive(event events.Event) error {
	queries := translateEvent(event)

	for _, query := range queries {
		b.queries <- query
	}

	return nil
}

func translateEvent(event events.Event) []Query {
	switch reflect.TypeOf(event) {
	case reflect.TypeOf(events.MoveEvent{}):
		return []Query{
			TurnNumberQuery(event.GameId()),
		}
	case reflect.TypeOf(events.GameStartEvent{}):
		gameStart := event.(*events.GameStartEvent)
		return []Query{
			UserGamesQuery(gameStart.WhiteId),
			UserGamesQuery(gameStart.BlackId),
		}
	case reflect.TypeOf(events.GameEndEvent{}):
		gameEnd := event.(*events.GameStartEvent)
		return []Query{
			UserGamesQuery(gameEnd.WhiteId),
			UserGamesQuery(gameEnd.BlackId),
		}
	case reflect.TypeOf(events.DrawOfferEvent{}):
		return []Query{
			DrawOfferStateQuery(event.GameId()),
		}
	case reflect.TypeOf(events.DrawOfferResponseEvent{}):
		return []Query{
			DrawOfferStateQuery(event.GameId()),
		}
	default:
		return []Query{}
	}
}
