package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type drawOfferStateQuery struct {
	GameId game.Id

	Answered bool
	result   drawOfferState
}

type drawOfferState string

const (
	noDrawOffer    drawOfferState = "none"
	blackDrawOffer drawOfferState = "black"
	whiteDrawOffer drawOfferState = "white"
)

func (q *drawOfferStateQuery) hasResult() bool {
	return q.Answered
}

func (q *drawOfferStateQuery) getResult() interface{} {
	return q.result
}

func (q *drawOfferStateQuery) computeResult(queries SystemQueries) {
	offers := queries.getEvents().EventsOfTypeForGame(q.GameId, events.DrawOfferType)
	responses := queries.getEvents().EventsOfTypeForGame(q.GameId, events.DrawOfferResponseType)

	q.Answered = true
	if len(responses) == len(offers) {
		q.result = noDrawOffer
		return
	}

	lastOffer := offers[len(offers)-1]
	if lastOffer.Offerer == game.White {
		q.result = whiteDrawOffer
	} else {
		q.result = blackDrawOffer
	}
}

func (q *drawOfferStateQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *drawOfferStateQuery) hash() string {
	return fmt.Sprintf("drawoffer:%v", q.GameId)
}
