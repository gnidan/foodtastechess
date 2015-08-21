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

	newFEN := InitializeFEN()
	stateFromFEN := newFEN.ConvertToState()
	stateFromOG := InitializeState()
	assert := assert.New(s.T())
	assert.Equal(stateFromOG, stateFromFEN)

}
