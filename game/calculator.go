package game

type GameCalculator interface {
	StartingFEN() FEN
	AfterMove(initial FEN, move AlgebraicMove) FEN
	ValidMoves(state FEN) []AlgebraicMove
}

type GameCalculatorService struct {
}

func NewGameCalculator() GameCalculator {
	return new(GameCalculatorService)
}

func (s *GameCalculatorService) StartingFEN() FEN {
	return FEN("")
}

func (s *GameCalculatorService) AfterMove(initial FEN, move AlgebraicMove) FEN {
	return FEN("")
}

func (s *GameCalculatorService) ValidMoves(state FEN) []AlgebraicMove {
	return []AlgebraicMove{}
}
