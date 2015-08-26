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

	nextFEN := AfterMove(AlgebraicMove("Pe2-e4"), newFEN) //first move, white pawn moves out 2
	assert.Equal(FEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"), nextFEN)
	//test to make sure you can't move opponent pieces (black attempt to move white pawn in a2)
	state := nextFEN.ConvertToState()
	assert.Equal([]AlgebraicMove{}, state.ValidMovesAtPos(NewPosition(1, 2)))

	nextFEN = AfterMove(AlgebraicMove("Pc7-c5"), nextFEN) //black moves pawn out 2
	assert.Equal(FEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Ng1-f3"), nextFEN) //white moves knight
	assert.Equal(FEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Pd7-d5"), nextFEN) //black moves pawn out 2
	assert.Equal(FEN("rnbqkbnr/pp2pppp/8/2pp4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq d6 0 3"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Pe4xd5"), nextFEN) //white pawn captures black pawn
	assert.Equal(FEN("rnbqkbnr/pp2pppp/8/2pP4/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 3"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Pe7-e5"), nextFEN) //black move pawn out 2

	//en passant move test
	nextFEN = AfterMove(AlgebraicMove("Pd5xe6.ep"), nextFEN) //white performs en passant
	assert.Equal(FEN("rnbqkbnr/pp3ppp/4P3/2p5/8/5N2/PPPP1PPP/RNBQKB1R b KQkq - 0 4"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Pd7xe6"), nextFEN) //black captures white pawn
	nextFEN = AfterMove(AlgebraicMove("Pg2-g3"), nextFEN) //white moves pawn out 1

	//castling availability test
	nextFEN = AfterMove(AlgebraicMove("Kd8-d7"), nextFEN) //black  moves king out (black no longer can castle)
	assert.Equal(FEN("rnb1kbnr/pp1k1ppp/4p3/2p5/8/5NP1/PPPP1P1P/RNBQKB1R w KQ - 1 6"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Bf1-g2"), nextFEN) //white moves bishop out 1
	nextFEN = AfterMove(AlgebraicMove("Qd8-h4"), nextFEN) //black moves queen out to h4

	//white performs castle kingside
	nextFEN = AfterMove(AlgebraicMove("0-0"), nextFEN) //white executes kingside castle
	assert.Equal(FEN("rnb1kbnr/pp1k1ppp/4p3/2p5/7q/5NP1/PPPP1PBP/RNBQ1RK1 b - - 1 7"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Kd7-e8"), nextFEN) //black moves king back to e8
	nextFEN = AfterMove(AlgebraicMove("Nf3-g5"), nextFEN) //white moves knight to g5
	nextFEN = AfterMove(AlgebraicMove("Nb8-a6"), nextFEN) //black moves knight to a6
	nextFEN = AfterMove(AlgebraicMove("Pa2-a3"), nextFEN) //white moves pawn to a3
	nextFEN = AfterMove(AlgebraicMove("Nc8-d7"), nextFEN) //black moves bishop to d7
	nextFEN = AfterMove(AlgebraicMove("Pa3-a4"), nextFEN) //white moves pawn to a4

	//making sure black still isnt allowed to castle, even though pieces are back in orig position
	assert.Equal(FEN("r3kbnr/pp1n1ppp/n3p3/2p3N1/P6q/6P1/1PPP1PBP/RNBQ1RK1 b - - 0 10"), nextFEN)
	state = nextFEN.ConvertToState()
	assert.Equal([]AlgebraicMove{AlgebraicMove("Ke8-d8"), AlgebraicMove("Ke8-e7")}, state.ValidMovesAtPos(NewPosition(5, 8)))

	nextFEN = AfterMove(AlgebraicMove("Pe6-e5"), nextFEN) //black moves pawn to e5
	nextFEN = AfterMove(AlgebraicMove("Pg3-g4"), nextFEN) //white moves pawn to g4
	nextFEN = AfterMove(AlgebraicMove("Pe5-e4"), nextFEN) //black moves pawn to e4
	nextFEN = AfterMove(AlgebraicMove("Bg2-h3"), nextFEN) //white moves bishop to h3
	nextFEN = AfterMove(AlgebraicMove("Qh4-g3"), nextFEN) //black moves queen to h3

	//not all possible moves are valid, because some would make white put self into check
	assert.NotEqual(AllPossibleMoves(nextFEN), AllValidMoves(nextFEN))

	nextFEN = AfterMove(AlgebraicMove("Ph2xg3"), nextFEN)   //white pawn captures black queen
	nextFEN = AfterMove(AlgebraicMove("Pe4-e3"), nextFEN)   //black pawn moves forward
	nextFEN = AfterMove(AlgebraicMove("Pd2-d4"), nextFEN)   //white pawn moves out 2
	nextFEN = AfterMove(AlgebraicMove("Pe3-e2"), nextFEN)   //black pawn moves forward
	nextFEN = AfterMove(AlgebraicMove("Ng5-e6"), nextFEN)   //white knight moves
	nextFEN = AfterMove(AlgebraicMove("Pe2xf1=Q"), nextFEN) //black pawn promotion to queen w/capture
	assert.Equal(FEN("r3kbnr/pp1n1ppp/n3N3/2p5/P2P2P1/6PB/1PP2P2/RNBQ1qK1 w - - 0 16"), nextFEN)

	nextFEN = AfterMove(AlgebraicMove("Kg1-h1"), nextFEN) //white king moves
	nextFEN = AfterMove(AlgebraicMove("Ph7-h5+"), nextFEN) //black pawn moves
	nextFEN = AfterMove(AlgebraicMove("Kh1-h2"), nextFEN) //white king moves
	nextFEN = AfterMove(AlgebraicMove("Ph5xg4"), nextFEN) //black pawn captures pawn
	nextFEN = AfterMove(AlgebraicMove("Ra1-a2"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Pg4xh3"), nextFEN) //black pawn captures bishop
	nextFEN = AfterMove(AlgebraicMove("Ra2-a3"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rh8-h4"), nextFEN) //black rook moves
	nextFEN = AfterMove(AlgebraicMove("Ra3-a2"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rh4xd4"), nextFEN) //black rook captures pawn
	nextFEN = AfterMove(AlgebraicMove("Ra2-a3"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rd4xd1"), nextFEN) //black captures queen
	nextFEN = AfterMove(AlgebraicMove("Ra3-a2"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rd1xc1"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Ra2-a3"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rc1xb1"), nextFEN) //black rook moves
	nextFEN = AfterMove(AlgebraicMove("Ra3-a2"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rb1x-b2"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Ra2-a3"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rb2xc2"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Ra3-a2"), nextFEN) //white rook moves
	nextFEN = AfterMove(AlgebraicMove("Rc2xa2"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Ne6-d4"), nextFEN) //white knight moves
	nextFEN = AfterMove(AlgebraicMove("Ra2xa4"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Pf2-f3"), nextFEN) //white pawn moves
	nextFEN = AfterMove(AlgebraicMove("Ra4xd4"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Pf2-f3"), nextFEN) //white pawn moves
	nextFEN = AfterMove(AlgebraicMove("Rd4-g4"), nextFEN) //black rook 
	nextFEN = AfterMove(AlgebraicMove("Pf3-f4"), nextFEN) //white pawn moves
	nextFEN = AfterMove(AlgebraicMove("Pb7-b6"), nextFEN) //black pawn
	nextFEN = AfterMove(AlgebraicMove("Pf4-f5"), nextFEN) //white pawn
	nextFEN = AfterMove(AlgebraicMove("Pb6-b5"), nextFEN) //black rook captures
	nextFEN = AfterMove(AlgebraicMove("Pf5-f6"), nextFEN) //white pawn 
	nextFEN = AfterMove(AlgebraicMove("Pg7xf6S"), nextFEN) //black rook captures, causes stalemate
	
	assert.Equal(FEN("r3kbn1/p2n1p2/n4p2/1pp5/6r1/6Pp/7K/5q2 w - - 0 33"), nextFEN)
	assert.Equal([]AlgebraicMove{}, AllValidMoves(nextFEN)) //no valid moves

}
