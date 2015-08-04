package main

type GameState struct{
	pieceMap [Position]Piece
}

type Position struct{
	file, rank int
}

func (s *GameState) PieceAtPosition(pos Position) piece Piece{
	piece, ok := pieceMap.Get(pos)
	if(ok == nul){
		//TODO error handling
	}else {
		return piece
	}
}
