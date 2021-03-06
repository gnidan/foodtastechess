package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/events"
	"foodtastechess/game"
)

type TurnNumberQueryTestSuite struct {
	QueryTestSuite
}

func (suite *TurnNumberQueryTestSuite) TestHasResult() {
	var (
		gameId              game.Id = 5
		hasResult, noResult *turnNumberQuery
	)

	hasResult = TurnNumberQuery(gameId).(*turnNumberQuery)
	hasResult.Result = 5

	noResult = TurnNumberQuery(gameId).(*turnNumberQuery)
	noResult.Result = -1

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *TurnNumberQueryTestSuite) TestDependentQueries() {
	var (
		gameId game.Id = 1
		query  *turnNumberQuery

		expectedDependents = []Query{}
	)

	query = TurnNumberQuery(gameId).(*turnNumberQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	assert.Equal(expectedDependents, actualDependents)
}

func (suite *TurnNumberQueryTestSuite) TestComputeResult() {
	var (
		gameId game.Id = 1
		query  *turnNumberQuery
	)

	query = TurnNumberQuery(gameId).(*turnNumberQuery)

	assert := assert.New(suite.T())

	// case 1

	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.MoveType).
		Return([]events.Event{}).
		Once()

	query.computeResult(suite.mockSystemQueries)

	assert.Equal(game.TurnNumber(0), query.Result)

	// case 2

	fakeMoves := []events.Event{
		events.NewMoveEvent(gameId, 1, ""),
		events.NewMoveEvent(gameId, 2, ""),
		events.NewMoveEvent(gameId, 3, ""),
	}
	suite.mockEvents.
		On("EventsOfTypeForGame", gameId, events.MoveType).
		Return(fakeMoves).
		Once()

	query.computeResult(suite.mockSystemQueries)
	assert.Equal(game.TurnNumber(3), query.Result)
}

func TestTurnNumberQueryTestSuite(t *testing.T) {
	suite.Run(t, new(TurnNumberQueryTestSuite))
}
