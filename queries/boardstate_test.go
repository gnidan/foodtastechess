package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BoardStateQueryTestSuite struct {
	suite.Suite
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

func TestBoardStateQueryTestSuite(t *testing.T) {
	suite.Run(t, new(BoardStateQueryTestSuite))
}
