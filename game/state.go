package game

type GameState struct {
	pieceMap map[Position]Piece
}

type Position struct {
	file, rank int
}

func NewGameState(pieceMap map[Position]Piece) GameState {
	gameState := GameState{}
	gameState.pieceMap = pieceMap
	return gameState
}

func NewPosition(file, rank int) Position {
	position := Position{}
	position.file = file
	position.rank = rank
	return position
}

func (s *GameState) PieceAtPosition(pos Position) Piece {
	piece, ok := s.pieceMap[pos]
	if ok {
		//piece found
		return piece
	}
	//piece not found
	return nil
}

func InitializeBoard() GameState {
	gameState := GameState{}
	gameState.pieceMap = map[Position]Piece{
		//White
		NewPosition(1, 2): NewPawn(White),
		NewPosition(2, 2): NewPawn(White),
		NewPosition(3, 2): NewPawn(White),
		NewPosition(4, 2): NewPawn(White),
		NewPosition(5, 2): NewPawn(White),
		NewPosition(6, 2): NewPawn(White),
		NewPosition(7, 2): NewPawn(White),
		NewPosition(8, 2): NewPawn(White),
		NewPosition(1, 1): NewRook(White),
		NewPosition(2, 1): NewKnight(White),
		NewPosition(3, 1): NewBishop(White),
		NewPosition(4, 1): NewQueen(White),
		NewPosition(5, 1): NewKing(White),
		NewPosition(6, 1): NewBishop(White),
		NewPosition(7, 1): NewKnight(White),
		NewPosition(8, 1): NewRook(White),

		//Black
		NewPosition(1, 7): NewPawn(Black),
		NewPosition(2, 7): NewPawn(Black),
		NewPosition(3, 7): NewPawn(Black),
		NewPosition(4, 7): NewPawn(Black),
		NewPosition(5, 7): NewPawn(Black),
		NewPosition(6, 7): NewPawn(Black),
		NewPosition(7, 7): NewPawn(Black),
		NewPosition(8, 7): NewPawn(Black),
		NewPosition(1, 8): NewRook(Black),
		NewPosition(2, 8): NewKnight(Black),
		NewPosition(3, 8): NewBishop(Black),
		NewPosition(4, 8): NewQueen(Black),
		NewPosition(5, 8): NewKing(Black),
		NewPosition(6, 8): NewBishop(Black),
		NewPosition(7, 8): NewKnight(Black),
		NewPosition(8, 8): NewRook(Black),
	}
	return gameState
}
