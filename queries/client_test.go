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
		gameId   game.Id         = 1
		expected game.TurnNumber = 5
		query                    = TurnNumberQuery(gameId)
	)
	suite.mockSystemQueries.On("GetAnswer", query).Return(expected)

	gameInfo := suite.ClientQueries.GameInformation(gameId)
	assert.Equal(suite.T(), expected, gameInfo.TurnNumber)
}

func TestClientQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientQueriesTestSuite))
}
