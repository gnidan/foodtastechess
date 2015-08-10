package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/directory"
	"foodtastechess/game"
	"foodtastechess/logger"
)

type BoardStateQueryTestSuite struct {
	suite.Suite

	log                  *logging.Logger
	mockSystemQueries    *MockSystemQueries    `inject:"systemQueries"`
	mockGameStateManager *MockGameStateManager `inject:"gameStateManager"`
}

func (suite *BoardStateQueryTestSuite) SetupTest() {
	suite.log = logger.Log("boardstate_test")

	var (
		d                directory.Directory
		systemQueries    MockSystemQueries
		gameStateManager MockGameStateManager
	)

	d = directory.New()
	d.AddService("systemQueries", &systemQueries)
	d.AddService("gameStateManager", &gameStateManager)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.mockSystemQueries = &systemQueries
	suite.mockGameStateManager = &gameStateManager
}

func (suite *BoardStateQueryTestSuite) TestHasResult() {
	var (
		hasResult, noResult *boardStateAtTurnQuery
	)

	hasResult = new(boardStateAtTurnQuery)
	hasResult.gameId = 5
	hasResult.turnNumber = 5
	hasResult.result = "Be5"

	noResult = new(boardStateAtTurnQuery)
	noResult.gameId = 5
	noResult.turnNumber = 5

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *BoardStateQueryTestSuite) TestDependentQueries() {
	var (
		gameId     game.Id         = 1
		turnNumber game.TurnNumber = 5
		query      *boardStateAtTurnQuery

		expectedDependents = []Query{
			BoardAtTurnQuery(gameId, turnNumber-1),
			MoveAtTurnQuery(gameId, turnNumber),
		}
	)

	query = BoardAtTurnQuery(gameId, turnNumber).(*boardStateAtTurnQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	for _, expected := range expectedDependents {
		assert.Contains(actualDependents, expected)
	}
}

func (suite *BoardStateQueryTestSuite) TestDependentQueriesBaseCase() {
	var (
		gameId     game.Id         = 1
		turnNumber game.TurnNumber = 0
		query      *boardStateAtTurnQuery

		expectedDependents = []Query{}
	)

	query = BoardAtTurnQuery(gameId, turnNumber).(*boardStateAtTurnQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	assert.Equal(expectedDependents, actualDependents)
}

func (suite *BoardStateQueryTestSuite) TestComputeResult() {
}

func TestBoardStateQueryTestSuite(t *testing.T) {
	suite.Run(t, new(BoardStateQueryTestSuite))
}
