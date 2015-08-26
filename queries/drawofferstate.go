package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type drawOfferStateQuery struct {
	GameId game.Id

	Answered bool
	Result   game.Color

	// Compose a queryRecord
	queryRecord `bson:",inline"`
}

func (q *drawOfferStateQuery) hasResult() bool {
	return q.Answered
}

func (q *drawOfferStateQuery) getResult() interface{} {
	return q.Result
}

func (q *drawOfferStateQuery) computeResult(queries SystemQueries) {
	offers := queries.getEvents().EventsOfTypeForGame(q.GameId, events.DrawOfferType)
	responses := queries.getEvents().EventsOfTypeForGame(q.GameId, events.DrawOfferResponseType)

	q.Answered = true
	if len(responses) == len(offers) {
		q.Result = game.NoOne
		return
	}

	lastOffer := offers[len(offers)-1]
	if lastOffer.Offerer == game.White {
		q.Result = game.White
	} else {
		q.Result = game.Black
	}
}

func (q *drawOfferStateQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *drawOfferStateQuery) hash() string {
	return fmt.Sprintf("drawoffer:%v", q.GameId)
}
