package game

type GameState struct{
	pieceMap map[Position]Piece
}

type Position struct{
	file, rank int
}

func (s *GameState) PieceAtPosition(pos Position) Piece{
	piece, ok := s.pieceMap[pos]
	if( ok ){
		//piece found
		return piece
	}
	//piece not found
	return nil
}
