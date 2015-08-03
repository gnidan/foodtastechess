package queries

import (
	"github.com/facebookgo/inject"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
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
	ClientQueries ClientQueries `inject:"clientQueries"`
}

type MockSystemQueries struct {
}

func (m *MockSystemQueries) GetAnswer(query Query) Answer {
	return 5
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
}

func (suite *ClientQueriesTestSuite) TestExample() {
	var expected game.TurnNumber
	expected = 5
	gameInfo := suite.ClientQueries.GameInformation(1)
	assert.Equal(suite.T(), expected, gameInfo.TurnNumber)
}

func TestClientQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientQueriesTestSuite))
}
