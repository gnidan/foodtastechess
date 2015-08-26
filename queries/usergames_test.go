package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/users"
)

type UserGamesQueryTestSuite struct {
	QueryTestSuite
}

func (suite *UserGamesQueryTestSuite) TestHasResult() {
	var (
		playerId            users.Id = "chloe"
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
		playerId users.Id = "bob"
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
		playerId      users.Id = "bob"
		query         *userGamesQuery
		activeGames   []game.Id = []game.Id{5, 6, 7}
		finishedGames []game.Id = []game.Id{1, 2, 3, 4}
		gameCreates   []events.Event
		gameStarts    []events.Event
		gameEnds      []events.Event
	)

	for _, id := range activeGames {
		gameCreates = append(gameCreates, events.NewGameCreateEvent(id, playerId, playerId))
		gameStarts = append(gameStarts, events.NewGameStartEvent(id, playerId, playerId))
	}
	for _, id := range finishedGames {
		gameCreates = append(gameCreates, events.NewGameCreateEvent(id, playerId, playerId))
		gameStarts = append(gameStarts, events.NewGameStartEvent(id, playerId, playerId))
		gameEnds = append(gameEnds, events.NewGameEndEvent(id, game.GameEndCheckmate, game.Black, playerId, playerId))
	}

	suite.mockEvents.
		On("EventsOfTypeForPlayer", playerId, events.GameCreateType).
		Return(gameCreates)

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
		assert.Contains(query.Result, id)
	}

	for _, id := range finishedGames {
		assert.Contains(query.Result, id)
	}

	assert.Equal(len(activeGames)+len(finishedGames), len(query.Result))

	assert.Equal(true, query.hasResult())
}

// Entrypoint
func TestUserGamesQueryTestSuite(t *testing.T) {
	suite.Run(t, new(UserGamesQueryTestSuite))
}
