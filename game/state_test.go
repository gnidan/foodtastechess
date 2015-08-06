package game

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GameStateTestSuite struct {
	suite.Suite
}

func TestGameStateTestSuite(t *testing.T) {
	suite.Run(t, new(GameStateTestSuite))
}

func (s *GameStateTestSuite) TestMapConstructor() {
	pieceMap := map[Position]Piece{
		NewPosition(2, 2): NewPawn(Black),
	}
	state := NewGameState(pieceMap)
	assert := assert.New(s.T())
	//assert.Equal(expected, actual)
	//assert.Equal(1,0) //force fail
	assert.Equal(state.PieceAtPosition(Position{2, 2}), NewPawn(Black))
}
