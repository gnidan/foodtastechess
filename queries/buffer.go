package queries

import (
	"github.com/op/go-logging"

	"foodtastechess/events"
	"foodtastechess/logger"
)

type QueryBuffer struct {
	log           *logging.Logger
	Queries       chan Query
	SystemQueries SystemQueries `inject:"systemQueries"`
	stopChan      chan bool
}

func NewQueryBuffer() events.EventSubscriber {
	buffer := new(QueryBuffer)
	buffer.log = logger.Log("querybuffer")
	buffer.Queries = make(chan Query, 100)
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
		case <-b.Queries:
			b.log.Info("Got query")
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
	return nil
}

func translateEvent(event events.Event) []Query {
	return []Query{}
}
