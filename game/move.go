package game

func (s *GameState) ValidMoves(pos Position) []Position{
	validMoves := [...]{}
	piece := s.PieceAtPosition(pos)
	if piece == nil{ //no piece at position
		return nil
	}
	moves := piece.Moves()
	for _, move := range moves {
		validMoves.append(move.Translate(pos, s))
	}
	return validMoves
}

//BoundMove - moves max of one position offset
//UnboundMove - moves until hit own piece, border, or opponent + 1
//AdvancingMove - Unique to Pawns: like bound move, but no capture allowed
//CapturingMove - Unique to Pawns: single capturing diagonal move forward
//EnPassantMove - Unique to Pawns: en passant
//CastlingMove - Unique to Rook, King (need for others???) - Castling

type BoundMove struct{
	rankOffset, fileOffset int
}

type UnboundMove struct{
	rankOffset, fileOffset int
}

type AdvancingMove struct{
	rankDelta int //white will move +1, black -1
}

type CapturingMove struct{
	rankOffset, fileOffset int
}

type EnPassantMove struct{
	rankOffset, fileOffset int
}

type CastlingMove struct{
	//TODO
}



func (m *BoundMove) Translate(pos Position, s GameState) []Position{
    newPos := Position{
		rank: pos.rank + m.rankOffset,
		file: pos.file + m.fileOffset,
		}
	//if there is no piece there, or piece is other color
    if s.PieceAtPosition(newPos) == nil
		|| s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
        return newPos
    }
	return nil
}

func (m *UnboundMove) Translate(pos Position, s GameState) []Position{
    newPos := new Position(pos.rank + m.rankOffset, pos.file + m.fileOffset)
    if s.PieceAtPosition(newPos) == nil {
        return [newPos].append(m.Translate(newPos, s))
    }else if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		return newPos
	}
	return nil
}

func (m *AdvancingMoveMove) Translate(pos Position, s GameState) []Position{
    newPos := new Position(pos.rank + m.rankDelta, pos.file)
    if s.PieceAtPosition(newPos) == nil {
        return newPos
    }
	return nil
}

func (m *CapturingMove) Translate(pos Position, s GameState) []Position{
    newPos := new Position(pos.rank + m.rankOffset, pos.file + m.fileOffset)
    if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
        return newPos
    }
	return nil
}

func (m *EnPassantMove) Translate(pos Position, s GameState) []Position{
    newPos := new Position(pos.rank + m.rankOffset, pos.file + m.fileOffset)
	oppPosition := new Position(pos.rank, pos.file + m.fileOffset)
    if s.PieceAtPosition(newPos) == nil 
		&& s.PieceAtPosition(oppPosition).Color() != s.PieceAtPosition(pos).Color() {
        //Other checks for valid en passant
		if {
			return newPos
		}
    }
    return nil
}

func (m *CastlingMove) Translate(pos Position, s GameState) []Position{
    //TODO
	return nil
}

