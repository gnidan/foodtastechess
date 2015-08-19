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

	/*
		pieceMap := map[Position]Piece{
			//NewPosition(2, 2): NewPawn(White),
			//NewPosition(3, 2): NewPawn(Black),
			NewPosition(1, 1): NewRook(White),
			NewPosition(8, 1): NewRook(White),
			NewPosition(5, 1): NewKing(White),
		}
		state := NewGameState(pieceMap)
		validMoves := state.ValidMoves(NewPosition(5, 1))


		//assert.Equal(expected, actual)
		//assert.Equal(0,1) //force fail

		assert.Equal(1, validMoves)
	*/

	newState := InitializeBoard()
	assert.Equal(newState.ConvertToFEN().ConvertToState().ConvertToFEN(), newState.ConvertToFEN())

}
