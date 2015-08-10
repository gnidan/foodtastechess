package queries

import (
	"foodtastechess/game"
)

type drawOfferStateQuery struct {
	gameId game.Id

	result struct {
		hasOffer bool
		offerer  game.Color
	}
}

func (q *drawOfferStateQuery) hasResult() bool {
	return true
}

func (q *drawOfferStateQuery) computeResult(map[Query]Query) {
}

func (q *drawOfferStateQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *drawOfferStateQuery) isExpired(now interface{}) bool {
	return false
}

func (q *drawOfferStateQuery) getExpiration(now interface{}) interface{} {
	return nil
}
