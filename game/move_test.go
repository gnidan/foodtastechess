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

	pieceMap := map[Position]Piece{
		//NewPosition(2, 2): NewPawn(White),
		//NewPosition(3, 2): NewPawn(Black),
		NewPosition(1, 1): NewRook(White),
		NewPosition(8, 1): NewRook(White),
		NewPosition(5, 1): NewKing(White),
	}
	state := NewGameState(pieceMap)

	//state := InitializeBoard()

	validMoves := state.ValidMoves(NewPosition(5, 1))

	assert := assert.New(s.T())
	//assert.Equal(expected, actual)
	//assert.Equal(0,1) //force fail
	assert.Equal(1, validMoves)
}
