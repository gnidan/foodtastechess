package game

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MoveTestSuite struct {
	suite.Suite
}

func TestMoveTestSuite(t *testing.T) {
	suite.Run(t, new(MoveTestSuite))
}

func (s *MoveTestSuite) TestMoveConstructor() {

	assert := assert.New(s.T())

	newState := InitializeState()
	assert.Equal(newState.ConvertToFEN().ConvertToState().ConvertToFEN(), newState.ConvertToFEN())
	newFEN := InitializeFEN()
	assert.Equal(newFEN, newState.ConvertToFEN())
	assert.Equal(newState, newFEN.ConvertToState())

	//after a first move
	move := AlgebraicMove("Pe2-e4")
	nextFEN := ExecuteMove(move, newFEN)
	assert.Equal(FEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"), nextFEN)

	//second move
	move = AlgebraicMove("Pc7-c5")
	nextFEN = ExecuteMove(move, nextFEN)
	assert.Equal(FEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"), nextFEN)

	move = AlgebraicMove("Ng1-f3")
	nextFEN = ExecuteMove(move, nextFEN)
	assert.Equal(FEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"), nextFEN)

	move = AlgebraicMove("Pd7-d5")
	nextFEN = ExecuteMove(move, nextFEN)
	assert.Equal(FEN("rnbqkbnr/pp2pppp/8/2pp4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq d6 0 3"), nextFEN)

	move = AlgebraicMove("Pe4xd5")
	nextFEN = ExecuteMove(move, nextFEN)
	assert.Equal(FEN("rnbqkbnr/pp2pppp/8/2pP4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 3"), nextFEN)

	//state := nextFEN.ConvertToState()
	//assert.Equal("1", state.ValidMoves(NewPosition(5,4)) )
}
