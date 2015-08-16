package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/directory"
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/logger"
)

type SystemQueriesTestSuite struct {
	suite.Suite

	log                *logging.Logger
	mockGameCalculator *MockGameCalculator
	mockEvents         *MockEventsService
	mockQueriesCache   *MockQueriesCache
	systemQueries      *SystemQueryService
}

func (suite *SystemQueriesTestSuite) SetupTest() {
	suite.log = logger.Log("system_test")

	var (
		d              directory.Directory
		gameCalculator MockGameCalculator
		events         MockEventsService
		queriesCache   MockQueriesCache
	)

	systemQueries := NewSystemQueryService().(*SystemQueryService)

	d = directory.New()
	d.AddService("systemQueries", systemQueries)
	d.AddService("gameCalculator", &gameCalculator)
	d.AddService("eventsService", &events)
	d.AddService("queriesCache", &queriesCache)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.systemQueries = systemQueries
	suite.mockGameCalculator = &gameCalculator
	suite.mockEvents = &events
	suite.mockQueriesCache = &queriesCache
}

func (suite *SystemQueriesTestSuite) TestLookup() {
	assert := assert.New(suite.T())

	var (
		gameId   game.Id         = 5
		expected game.TurnNumber = 10
		query    Query           = TurnNumberQuery(gameId)
	)

	suite.mockQueriesCache.
		On("Get", query).
		Return(true).
		Run(
		func(args mock.Arguments) {
			partial := args.Get(0).(*turnNumberQuery)
			partial.result = expected
		})

	actual := suite.systemQueries.AnswerQuery(query).(game.TurnNumber)
	assert.Equal(expected, actual)
}

func (suite *SystemQueriesTestSuite) TestCompute() {
	assert := assert.New(suite.T())

	var (
		gameId game.Id = 5
		query  Query   = TurnNumberQuery(gameId)
	)

	suite.mockQueriesCache.
		On("Get", query).
		Return(false)

	fakeMoves := []events.Event{
		events.NewMoveEvent(gameId, 1, ""),
		events.NewMoveEvent(gameId, 2, ""),
		events.NewMoveEvent(gameId, 3, ""),
	}
	expected := game.TurnNumber(len(fakeMoves))
	suite.mockEvents.On("EventsOfTypeForGame", gameId, events.MoveType).Return(fakeMoves).Once()

	suite.mockQueriesCache.On("Store", query).Return().Once()

	actual := suite.systemQueries.AnswerQuery(query).(game.TurnNumber)
	assert.Equal(expected, actual)
}

func (suite *SystemQueriesTestSuite) TestDependentQueries() {
	assert := assert.New(suite.T())

	var (
		gameId                   game.Id            = 5
		turnNumber               game.TurnNumber    = 11
		previousBoardState       game.FEN           = "previous!"
		previousBoardAtTurnQuery Query              = BoardAtTurnQuery(gameId, turnNumber-1)
		lastMove                 game.AlgebraicMove = "move!"
		lastMoveQuery            Query              = MoveAtTurnQuery(gameId, turnNumber)
		expectedState            game.FEN           = "current!"
		query                    Query              = BoardAtTurnQuery(gameId, turnNumber)
	)

	suite.mockQueriesCache.
		On("Get", previousBoardAtTurnQuery).
		Return(true).
		Run(
		func(args mock.Arguments) {
			partial := args.Get(0).(*boardStateAtTurnQuery)
			partial.result = previousBoardState
		})

	suite.mockQueriesCache.
		On("Get", lastMoveQuery).
		Return(true).
		Run(
		func(args mock.Arguments) {
			partial := args.Get(0).(*moveAtTurnQuery)
			partial.result = lastMove
		})

	suite.mockQueriesCache.
		On("Get", query).
		Return(false)

	suite.mockGameCalculator.
		On("AfterMove", previousBoardState, lastMove).
		Return(expectedState)

	suite.mockQueriesCache.On("Store", query).Return().Once()

	actualState := suite.systemQueries.AnswerQuery(query).(game.FEN)
	assert.Equal(expectedState, actualState)
}

func (suite *SystemQueriesTestSuite) TestQueryStorage() {
	var (
		gameId game.Id = 5
		query  Query   = TurnNumberQuery(gameId)
	)

	suite.mockQueriesCache.
		On("Get", query).
		Return(false)

	fakeMoves := []events.Event{
		events.NewMoveEvent(gameId, 1, ""),
		events.NewMoveEvent(gameId, 2, ""),
		events.NewMoveEvent(gameId, 3, ""),
	}
	suite.mockEvents.On("EventsOfTypeForGame", gameId, events.MoveType).Return(fakeMoves).Once()

	suite.mockQueriesCache.On("Store", query).Return().Once()

	suite.systemQueries.AnswerQuery(query)
	suite.mockQueriesCache.AssertCalled(suite.T(), "Store", query)
}

func TestSystemQueries(t *testing.T) {
	suite.Run(t, new(SystemQueriesTestSuite))
}
