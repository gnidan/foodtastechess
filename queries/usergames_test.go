package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/events"
	"foodtastechess/game"
)

type UserGamesQueryTestSuite struct {
	QueryTestSuite
}

func (suite *UserGamesQueryTestSuite) TestHasResult() {
	var (
		playerId            string = "5xvf"
		hasResult, noResult *userGamesQuery
	)

	hasResult = UserGamesQuery(playerId).(*userGamesQuery)
	hasResult.Answered = true

	noResult = UserGamesQuery(playerId).(*userGamesQuery)
	noResult.Answered = false

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *UserGamesQueryTestSuite) TestDependentQueries() {
	var (
		playerId string = "5xvf"
		query    *userGamesQuery

		expectedDependents = []Query{}
	)

	query = UserGamesQuery(playerId).(*userGamesQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	assert.Equal(expectedDependents, actualDependents)
}

func (suite *UserGamesQueryTestSuite) TestComputeResult() {
	var (
		playerId      string = "bob"
		query         *userGamesQuery
		activeGames   []game.Id = []game.Id{5, 6, 7}
		finishedGames []game.Id = []game.Id{1, 2, 3, 4}
		gameStarts    []events.Event
		gameEnds      []events.Event
	)

	for _, id := range activeGames {
		gameStarts = append(gameStarts, events.NewGameStartEvent(id))
	}
	for _, id := range finishedGames {
		gameStarts = append(gameStarts, events.NewGameStartEvent(id))
		gameEnds = append(gameEnds, events.NewGameEndEvent(id))
	}

	suite.mockEvents.
		On("EventsOfTypeForPlayer", playerId, events.GameStartType).
		Return(gameStarts)

	suite.mockEvents.
		On("EventsOfTypeForPlayer", playerId, events.GameEndType).
		Return(gameEnds)

	query = UserGamesQuery(playerId).(*userGamesQuery)

	assert := assert.New(suite.T())

	assert.Equal(false, query.hasResult())

	query.computeResult(suite.mockSystemQueries)

	for _, id := range activeGames {
		assert.Contains(query.result, id)
	}

	assert.Equal(len(activeGames), len(query.result))

	assert.Equal(true, query.hasResult())
}

func TestUserGamesQueryTestSuite(t *testing.T) {
	suite.Run(t, new(UserGamesQueryTestSuite))
}
