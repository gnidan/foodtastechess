package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/game"
)

type ValidMovesQueryTestSuite struct {
	QueryTestSuite
}

func (suite *ValidMovesQueryTestSuite) TestHasResult() {
	var (
		hasResult, noResult *validMovesAtTurnQuery
	)

	hasResult = new(validMovesAtTurnQuery)
	hasResult.GameId = 5
	hasResult.TurnNumber = 11
	hasResult.Result = []game.AlgebraicMove{"Qe7"}
	hasResult.Answered = true

	noResult = new(validMovesAtTurnQuery)
	noResult.GameId = 2
	noResult.TurnNumber = 6
	noResult.Result = []game.AlgebraicMove{"Ra1"}

	assert := assert.New(suite.T())
	assert.Equal(true, hasResult.hasResult())
	assert.Equal(false, noResult.hasResult())
}

func (suite *ValidMovesQueryTestSuite) TestDependentQueries() {
	var (
		gameId     game.Id         = 1
		turnNumber game.TurnNumber = 8
		query      Query

		expectedDependents = []Query{
			BoardAtTurnQuery(gameId, turnNumber),
		}
	)

	query = ValidMovesAtTurnQuery(gameId, turnNumber)

	actualDependents := query.getDependentQueries()

	assert := assert.New(suite.T())
	for _, expected := range expectedDependents {
		assert.Contains(actualDependents, expected)
	}
}

func (suite *ValidMovesQueryTestSuite) TestComputeResult() {
	var (
		gameId     game.Id         = 1
		turnNumber game.TurnNumber = 5

		expectedState game.FEN = "gonnawin!"

		expectedValidMoves = []game.AlgebraicMove{
			"goodmove",
			"badmove",
			"bestmove",
		}

		boardStateQ Query = &boardStateAtTurnQuery{
			GameId:     gameId,
			TurnNumber: 1,
			Result:     expectedState,
		}

		validMovesQ *validMovesAtTurnQuery = ValidMovesAtTurnQuery(gameId, turnNumber).(*validMovesAtTurnQuery)
	)

	assert := assert.New(suite.T())

	suite.mockSystemQueries.
		On("getDependentQueryLookup", validMovesQ).
		Return(NewQueryLookup(boardStateQ))

	suite.mockGameCalculator.
		On("ValidMoves", expectedState).
		Return(expectedValidMoves)

	validMovesQ.computeResult(suite.mockSystemQueries)
	assert.Equal(expectedValidMoves, validMovesQ.Result)
}

// Entrypoint
func TestValidMovesQueryTestSuite(t *testing.T) {
	suite.Run(t, new(ValidMovesQueryTestSuite))
}
