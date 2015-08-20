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
	hasResult.Result = []game.MoveRecord{
		game.MoveRecord{"horsey", "checkmate"},
	}
	hasResult.Answered = true

	noResult = new(validMovesAtTurnQuery)
	noResult.GameId = 2
	noResult.TurnNumber = 6
	noResult.Result = []game.MoveRecord{}

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

		initialState game.FEN = "a daring battle of wits"

		fakeMoves = []game.AlgebraicMove{
			"goodmove",
			"badmove",
			"bestmove",
		}

		fakeOutcomes = []game.FEN{
			"check!",
			"ball in hand ?!",
			"checkmate",
		}

		boardStateQ Query = &boardStateAtTurnQuery{
			GameId:     gameId,
			TurnNumber: turnNumber,
			Result:     initialState,
		}

		// we'll calculate the pairs below
		expectedResult []game.MoveRecord

		validMovesQ *validMovesAtTurnQuery = ValidMovesAtTurnQuery(gameId, turnNumber).(*validMovesAtTurnQuery)
	)

	assert := assert.New(suite.T())

	suite.mockSystemQueries.
		On("getDependentQueryLookup", validMovesQ).
		Return(NewQueryLookup(boardStateQ)).
		Once()

	suite.mockGameCalculator.
		On("ValidMoves", initialState).
		Return(fakeMoves).
		Once()

	expectedResult = []game.MoveRecord{}
	for i, move := range fakeMoves {
		outcome := fakeOutcomes[i]

		expectedResult = append(expectedResult, game.MoveRecord{
			Move:                move,
			ResultingBoardState: outcome,
		})

		suite.mockGameCalculator.
			On("AfterMove", initialState, move).
			Return(outcome).
			Once()
	}

	validMovesQ.computeResult(suite.mockSystemQueries)
	assert.Equal(expectedResult, validMovesQ.Result)
}

// Entrypoint
func TestValidMovesQueryTestSuite(t *testing.T) {
	suite.Run(t, new(ValidMovesQueryTestSuite))
}
