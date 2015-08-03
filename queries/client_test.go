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

type ClientQueriesTestSuite struct {
	suite.Suite
	mockSystemQueries *MockSystemQueries
	ClientQueries     ClientQueries `inject:"clientQueries"`
}

type MockSystemQueries struct {
	mock.Mock
}

func (m *MockSystemQueries) GetAnswer(query Query) Answer {
	args := m.Called(query)
	return args.Get(0).(Answer)
}

func (suite *ClientQueriesTestSuite) SetupTest() {
	var (
		systemQueries MockSystemQueries
		service       ClientQueryService
	)

	var g inject.Graph
	err := g.Provide(
		&inject.Object{Value: suite},
		&inject.Object{Name: "systemQueries", Value: &systemQueries},
		&inject.Object{Name: "clientQueries", Value: &service},
	)
	if err != nil {
		log.Fatalf("Could not provide objects to graph, %v", err)
	}

	if err = g.Populate(); err != nil {
		log.Fatalf("Could not populate graph %v", err)
	}

	suite.mockSystemQueries = &systemQueries
}

func (suite *ClientQueriesTestSuite) TestExample() {
	var (
		gameId             game.Id         = 1
		expectedTurnNumber game.TurnNumber = 5
		expectedBoardState game.FEN        = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
		turnNumberQuery    Query           = TurnNumberQuery(gameId)
		boardStateQuery    Query           = BoardAtTurnQuery(gameId, expectedTurnNumber)
	)
	suite.mockSystemQueries.On("GetAnswer", turnNumberQuery).Return(expectedTurnNumber)
	suite.mockSystemQueries.On("GetAnswer", boardStateQuery).Return(expectedBoardState)

	gameInfo := suite.ClientQueries.GameInformation(gameId)
	assert.Equal(suite.T(), expectedTurnNumber, gameInfo.TurnNumber)
	assert.Equal(suite.T(), expectedBoardState, gameInfo.BoardState)
}

func TestClientQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientQueriesTestSuite))
}
