package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/events"
	"foodtastechess/game"
)

type MoveAtTurnQueryTestSuite struct {
	QueryTestSuite
}

func (suite *MoveAtTurnQueryTestSuite) TestHasResult() {
	var (
		gameId              game.Id         = 5
		turnNumber          game.TurnNumber = 9
		hasResult, noResult *moveAtTurnQuery
	)

	hasResult = MoveAtTurnQuery(gameId, turnNumber).(*moveAtTurnQuery)
	hasResult.Result = game.AlgebraicMove("Be5")
	hasResult.Answered = true

	noResult = MoveAtTurnQuery(gameId, turnNumber).(*moveAtTurnQuery)
	noResult.Result = game.AlgebraicMove("")
	noResult.Answered = false

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *MoveAtTurnQueryTestSuite) TestDependentQueries() {
	var (
		gameId     game.Id         = 1
		turnNumber game.TurnNumber = 5
		query      *moveAtTurnQuery

		expectedDependents = []Query{}
	)

	query = MoveAtTurnQuery(gameId, turnNumber).(*moveAtTurnQuery)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	assert.Equal(expectedDependents, actualDependents)
}

func (suite *MoveAtTurnQueryTestSuite) TestComputeResult() {
	var (
		gameId     game.Id            = 7
		turnNumber game.TurnNumber    = 9
		move       game.AlgebraicMove = "Na2"

		query *moveAtTurnQuery
		event events.Event
	)

	event = events.NewMoveEvent(gameId, turnNumber, move)
	query = MoveAtTurnQuery(gameId, turnNumber).(*moveAtTurnQuery)

	assert := assert.New(suite.T())
	suite.mockEvents.
		On("MoveEventForGameAtTurn", gameId, turnNumber).
		Return(event)

	query.computeResult(suite.mockSystemQueries)
	assert.Equal(true, query.hasResult())
	assert.Equal(move, query.Result)
}

func TestMoveAtTurnQueryTestSuite(t *testing.T) {
	suite.Run(t, new(MoveAtTurnQueryTestSuite))
}
