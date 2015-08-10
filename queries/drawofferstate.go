package queries

import (
	"fmt"

	"foodtastechess/events"
	"foodtastechess/game"
)

type drawOfferStateQuery struct {
	gameId game.Id

	answered bool
	result   drawOfferState
}

type drawOfferState string

const (
	noDrawOffer    drawOfferState = "none"
	blackDrawOffer drawOfferState = "black"
	whiteDrawOffer drawOfferState = "white"
)

func (q *drawOfferStateQuery) hasResult() bool {
	return q.answered
}

func (q *drawOfferStateQuery) computeResult(queries SystemQueries) {
	offers := queries.getEvents().EventsOfTypeForGame(q.gameId, events.DrawOfferType)
	responses := queries.getEvents().EventsOfTypeForGame(q.gameId, events.DrawOfferResponseType)

	q.answered = true
	if len(responses) == len(offers) {
		q.result = noDrawOffer
		return
	}

	lastOffer := offers[len(offers)-1].(events.DrawOfferEvent)
	if lastOffer.Color == game.White {
		q.result = whiteDrawOffer
	} else {
		q.result = blackDrawOffer
	}
}

func (q *drawOfferStateQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *drawOfferStateQuery) hash() string {
	return fmt.Sprintf("drawoffer:%v", q.gameId)
}
