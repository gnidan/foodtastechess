package game

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PieceTestSuite struct {
	suite.Suite
}

func TestPieceTestSuite(t *testing.T) {
	suite.Run(t, new(PieceTestSuite))
}

func (s *PieceTestSuite) TestPieceConstructor() {
	pawn := NewPawn(Black)
	assert := assert.New(s.T())
	assert.Equal(pawn.color, Black)
}
