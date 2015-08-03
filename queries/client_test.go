package queries

import (
	"github.com/facebookgo/inject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/game"
	"foodtastechess/logger"
)

var (
	log = logger.Log("client_test")
)

// ClientQueriesTestSuite is a collection of tests to ensure the correct
// behavior of the Client Query Service (CQS.) The tests utilize a mocked
// System Query Service (SQS) in order to test that the CQS is merely
// aggregating results from the system.
type ClientQueriesTestSuite struct {
	suite.Suite

	mockSystemQueries *MockSystemQueries
	clientQueries     ClientQueries
}

// MockSystemQueries is a mock that we're going to use as a
// SystemQueryInterface
type MockSystemQueries struct {
	mock.Mock
}

// GetAnswer records the call with Query and returns the pre-configured
// mock answer
func (m *MockSystemQueries) GetAnswer(query Query) Answer {
	args := m.Called(query)
	return args.Get(0).(Answer)
}

// SetupTest prepares the test suite for running by making a fake system
// query service, providing it to a real client query service (the one we
// are testing)
func (suite *ClientQueriesTestSuite) SetupTest() {
	var (
		g inject.Graph

		systemQueries MockSystemQueries
		clientQueries ClientQueryService
	)

	// Set up the graph with:
	//  - A real ClientQueryService (The one we are testing)
	//  - The mocked SystemQueries implementation
	if err := g.Provide(
		&inject.Object{Name: "clientQueries", Value: &clientQueries},
		&inject.Object{Name: "systemQueries", Value: &systemQueries},
	); err != nil {
		log.Fatalf("Could not provide values (%v)", err)
	}

	// Populate the graph so that clientQueries knows to use our mocked
	// systemQueries
	if err := g.Populate(); err != nil {
		log.Fatalf("Could not populate graph (%v)", err)
	}

	// Store references for use in tests
	suite.mockSystemQueries = &systemQueries
	suite.clientQueries = &clientQueries
}

// TestGameInformation tests the ClientQueries.GameInformation() method.
//
// GameInformation should query the SQS for the current turn number and
// the board state at that turn, and return a GameInformation struct
// with that information.
func (suite *ClientQueriesTestSuite) TestGameInformation() {
	var (
		// the game ID we'll be using
		gameId game.Id = 1

		// pretend it's this turn
		expectedTurnNumber game.TurnNumber = 5

		// in this board state
		expectedBoardState game.FEN = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"

		// expected query objects we're looking for
		turnNumberQuery Query = TurnNumberQuery(gameId)
		boardStateQuery Query = BoardAtTurnQuery(gameId, expectedTurnNumber)
	)

	// given our expected queries, return our respective expected results
	suite.mockSystemQueries.On("GetAnswer", turnNumberQuery).Return(expectedTurnNumber)
	suite.mockSystemQueries.On("GetAnswer", boardStateQuery).Return(expectedBoardState)

	// run the test call
	gameInfo := suite.clientQueries.GameInformation(gameId)

	// and expect that the game info we get back has the pretend values
	assert.Equal(suite.T(), expectedTurnNumber, gameInfo.TurnNumber)
	assert.Equal(suite.T(), expectedBoardState, gameInfo.BoardState)
}

func TestClientQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientQueriesTestSuite))
}
