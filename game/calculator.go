package game

type GameCalculator interface {
	StartingFEN() FEN
	AfterMove(initial FEN, move AlgebraicMove) FEN
}

type GameCalculatorService struct {
}
