package game

type GameCalculator interface {
	StartingFEN() FEN
	AfterMove(initial FEN, move AlgebraicMove) FEN
	ValidMoves(state FEN) []AlgebraicMove
}

type GameCalculatorService struct {
}
