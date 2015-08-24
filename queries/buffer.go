package queries

import (
	"github.com/op/go-logging"

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
			b.log.Info("QueryBuffer stopped")
			return
		}
	}
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
	switch event.Type {
	case events.MoveType:
		return []Query{
			TurnNumberQuery(event.GameId),
		}
	case events.GameCreateType:
		queries := []Query{
			GameQuery(event.GameId),
		}

		if event.WhiteId != "" {
			queries = append(queries, UserGamesQuery(event.WhiteId))
		}

		if event.BlackId != "" {
			queries = append(queries, UserGamesQuery(event.BlackId))
		}

		return queries
	case events.GameStartType:
		return []Query{
			GameQuery(event.GameId),
			GamePlayersQuery(event.GameId),
			UserGamesQuery(event.WhiteId),
			UserGamesQuery(event.BlackId),
		}
	case events.GameEndType:
		return []Query{
			GameQuery(event.GameId),
			UserGamesQuery(event.WhiteId),
			UserGamesQuery(event.BlackId),
		}
	case events.DrawOfferType:
		return []Query{
			DrawOfferStateQuery(event.GameId),
		}
	case events.DrawOfferResponseType:
		return []Query{
			DrawOfferStateQuery(event.GameId),
		}
	default:
		return []Query{}
	}
}
