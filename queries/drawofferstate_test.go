package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/events"
	"foodtastechess/game"
)

type DrawOfferStateQueryTestSuite struct {
	QueryTestSuite
}

func (suite *DrawOfferStateQueryTestSuite) TestHasResult() {
	var (
		gameId              game.Id = 5
		hasResult, noResult *drawOfferStateQuery
	)

	hasResult = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)
	hasResult.answered = true

	noResult = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)
	noResult.answered = false

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *DrawOfferStateQueryTestSuite) TestDependentQueries() {
	var (
		gameId game.Id = 1
		query  *drawOfferStateQuery

		expectedDependents = []Query{}
	)

	query = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	assert.Equal(expectedDependents, actualDependents)
}

func (suite *DrawOfferStateQueryTestSuite) TestComputeResult() {
	assert := assert.New(suite.T())

	var (
		gameId game.Id
		query  *drawOfferStateQuery
	)

	// Offer -> No Response
	gameId = 1
	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferType).
		Return([]events.Event{
		events.NewDrawOfferEvent(gameId, game.White),
	})
	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferResponseType).
		Return([]events.Event{})

	query = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)
	query.computeResult(suite.mockSystemQueries)
	assert.Equal(whiteDrawOffer, query.result)

	// Offer -> Accept
	gameId = 2
	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferType).
		Return([]events.Event{
		events.NewDrawOfferEvent(gameId, game.Black),
	})

	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferResponseType).
		Return([]events.Event{
		events.NewDrawOfferResponseEvent(gameId, true),
	})

	query = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)
	query.computeResult(suite.mockSystemQueries)
	assert.Equal(noDrawOffer, query.result)

	// Offer -> Reject -> New Offer
	gameId = 3
	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferType).
		Return([]events.Event{
		events.NewDrawOfferEvent(gameId, game.White),
		events.NewDrawOfferEvent(gameId, game.Black),
	})

	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.DrawOfferResponseType).
		Return([]events.Event{
		events.NewDrawOfferResponseEvent(gameId, false),
	})

	query = DrawOfferStateQuery(gameId).(*drawOfferStateQuery)
	query.computeResult(suite.mockSystemQueries)
	assert.Equal(blackDrawOffer, query.result)

}

func TestDrawOfferStateQueryTestSuite(t *testing.T) {
	suite.Run(t, new(DrawOfferStateQueryTestSuite))
}
