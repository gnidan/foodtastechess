package game

//BoundMove - moves max of one position offset
//UnboundMove - moves until hit own piece, border, or opponent + 1
//AdvancingMove - Unique to Pawns: like bound move, but no capture allowed
//CapturingMove - Unique to Pawns: single capturing diagonal move forward
//EnPassantMove - Unique to Pawns: en passant
//CastlingMove - Unique to Rook, King (need for others???) - Castling

type Move interface {
	Translate(pos Position, s *GameState) []Position
}

type BoundMove struct {
	fileOffset, rankOffset int
}

type UnboundMove struct {
	fileOffset, rankOffset int
}

type FirstPawnMove struct {
	rankDelta int
}

type AdvancingMove struct {
	rankDelta int //white will move +1, black -1
}

type CapturingMove struct {
	fileOffset, rankOffset int
}

type EnPassantMove struct {
	fileOffset, rankOffset int
}

type CastlingMove struct {
}


//Takes FEN and Algebraic, executes move and returns string
func ExecuteMove(string FEN string, stringAlgabraic string) string {
	
}


func (s *GameState) ValidMoves(pos Position) []Position {
	validMoves := []Position{} //make empty array of positions
	piece := s.PieceAtPosition(pos)
	if piece == nil { //no piece at position
		return nil
	}
	moves := piece.Moves()
	for _, move := range moves {
		validMoves = append(validMoves, move.Translate(pos, s)...)
	}
	return validMoves
}

func (m *BoundMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	//make sure position is within boundaries
	if newPos.file < 1 || newPos.file > 8 || newPos.rank < 1 || newPos.rank > 8 {
		return nil
	}

	//if there is no piece there, or piece is other color
	if s.PieceAtPosition(newPos) == nil || s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		return []Position{newPos}
	}
	return nil
}

/*
func (m *UnboundMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	//make sure position is within boundaries
	if newPos.file < 1 || newPos.file > 8 || newPos.rank < 1 || newPos.rank > 8 {
		return nil
	}
	//if no piece at position, it is valid
	if s.PieceAtPosition(newPos) == nil {
		return append([]Position{newPos}, m.Translate(newPos, s)...)
	}
	//there is a piece in next position
	if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		//if opponent is in position, it is valid, but end of recursion
		return []Position{newPos}
	}
	return nil
}
*/

func (m *UnboundMove) Translate(pos Position, s *GameState) []Position {

	posList := []Position{}

	nextPos := pos

	for {
		nextPos.file = nextPos.file + m.fileOffset
		nextPos.rank = nextPos.rank + m.rankOffset

		//make sure position is within boundaries
		if nextPos.file < 1 || nextPos.file > 8 || nextPos.rank < 1 || nextPos.rank > 8 {

			return posList
		} else if s.PieceAtPosition(nextPos) == nil {
			//if no piece at position, it is valid.  keep looping
			posList = append(posList, nextPos)
		} else if s.PieceAtPosition(nextPos).Color() != s.PieceAtPosition(pos).Color() {
			//there is a piece in next position

			//if opponent is in position, it is valid, but end of recursion
			return append(posList, nextPos)
		} else {
			//Piece of same color is in next pos.  end of recursion
			return posList
		}
	}

	return nil //should not hit this
}

func (m *FirstPawnMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file,
		rank: pos.rank + m.rankDelta,
	}
	midPos := Position{
		file: pos.file,
		rank: pos.rank + m.rankDelta/2,
	}
	//if both spaces in front of pawn are open, move is valid
	if s.PieceAtPosition(newPos) == nil && s.PieceAtPosition(midPos) == nil {
		//if rank is same as initial, move is valid
		pawn := s.PieceAtPosition(pos)
		if pawn.Color() == White && pos.rank == 2 {
			return []Position{newPos}
		}
		if pawn.Color() == Black && pos.rank == 7 {
			return []Position{newPos}
		}
	}
	return nil
}

func (m *AdvancingMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file,
		rank: pos.rank + m.rankDelta,
	}
	if s.PieceAtPosition(newPos) == nil {
		return []Position{newPos}
	}
	return nil
}

func (m *CapturingMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	//no opponent in position to be captured, invalid
	if s.PieceAtPosition(newPos) == nil {
		return nil
	}
	//piece in pos must be opponent
	if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		return []Position{newPos}
	}
	return nil
}

func (m *EnPassantMove) Translate(pos Position, s *GameState) []Position {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	oppPosition := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank,
	}

	//TODO Check if previous move was opponent moving pawn (to be captured) out 2 spaces

	//no opponent in place for en passant, or position to move in is occupied.  invalid move
	if s.PieceAtPosition(oppPosition) == nil || s.PieceAtPosition(newPos) != nil {
		return nil
	}
	//make sure piece in opponent pos is actually opponent
	if s.PieceAtPosition(oppPosition).Color() != s.PieceAtPosition(pos).Color() {
		//Make sure pawn is in correct rank to make en passant
		pawn := s.PieceAtPosition(pos)
		if pawn.Color() == White && pos.rank == 5 {
			return []Position{newPos}
		}
		if pawn.Color() == Black && pos.rank == 4 {
			return []Position{newPos}
		}
	}
	return nil
}

func (m *CastlingMove) Translate(pos Position, s *GameState) []Position {
	//TODO test if King or Rook has been moved previously

	posList := []Position{} //need to do this in array form, in case 2 castling moves possible

	king := s.PieceAtPosition(pos)
	//Piece must be in initial rank
	if king.Color() == White && pos.file == 5 && pos.rank == 1 {
		//White
		//make sure nothing in between rook and king
		if s.PieceAtPosition(NewPosition(2, 1)) == nil && s.PieceAtPosition(NewPosition(3, 1)) == nil && s.PieceAtPosition(NewPosition(4, 1)) == nil {
			//left side
			posList = append(posList, NewPosition(3, 1))
		}
		if s.PieceAtPosition(NewPosition(6, 1)) == nil && s.PieceAtPosition(NewPosition(7, 1)) == nil {
			//right side
			posList = append(posList, NewPosition(7, 1))
		}
	} else if king.Color() == Black && pos.file == 5 && pos.rank == 8 {
		//Black
		//make sure nothing in between rook and king
		if s.PieceAtPosition(NewPosition(2, 8)) == nil && s.PieceAtPosition(NewPosition(3, 8)) == nil && s.PieceAtPosition(NewPosition(4, 8)) == nil {
			//left side
			posList = append(posList, NewPosition(3, 8))
		}
		if s.PieceAtPosition(NewPosition(6, 8)) == nil && s.PieceAtPosition(NewPosition(7, 8)) == nil {
			//right side
			posList = append(posList, NewPosition(7, 8))
		}
	}
	return posList
}
