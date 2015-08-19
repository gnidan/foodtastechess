package game

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	//"strconv"
)

type GameStateTestSuite struct {
	suite.Suite
}

func TestGameStateTestSuite(t *testing.T) {
	suite.Run(t, new(GameStateTestSuite))
}

func (s *GameStateTestSuite) TestMapConstructor() {
	/*
		pieceMap := map[Position]Piece{}
		pieceMap[NewPosition(2, 2)] = NewPawn(Black)

		state := NewGameState(pieceMap)
		assert := assert.New(s.T())
		//assert.Equal(expected, actual)
		//assert.Equal(1,0) //force fail
		assert.Equal(state.PieceAtPosition(Position{2, 2}), NewPawn(Black))
	*/

	state := InitializeBoard()
	//pieceMap := map[Position]Piece{}
	//pieceMap[NewPosition(2, 2)] = NewPawn(Black)
	//state := NewGameState(pieceMap)
	newFEN := state.ConvertToFEN()

	newStateFromFEN := newFEN.ConvertToState()

	assert := assert.New(s.T())
	assert.Equal(state, newStateFromFEN)
	//assert.Equal(state, newFEN)

	/*
		stringFEN := "GBCDEFG"
		next := stringFEN[:1]

		temp,_ := strconv.Atoi(next)
		assert := assert.New(s.T())
		assert.Equal("...", temp)
	*/
}
