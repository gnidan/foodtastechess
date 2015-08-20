package game

import (
	"strconv"
)

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

//Takes AlgebraicMove and FEN, executes move and returns resulting FEN
func ExecuteMove(algebraicMove AlgebraicMove, prevFEN FEN) FEN {
	state := prevFEN.ConvertToState()
	state.enPassantTarget = ""
	state.halfmoveClock += 1
	color := state.activeColor
	newPieceMap := state.pieceMap
	stringAN := string(algebraicMove)

	if stringAN == "0-0" {
		// KINGSIDE CASTLING:
		castleRank := 0
		if color == White {
			castleRank = 1
		} else {
			castleRank = 8
		}
		newPieceMap[NewPosition(5, castleRank)] = nil
		newPieceMap[NewPosition(8, castleRank)] = nil
		newPieceMap[NewPosition(7, castleRank)] = NewKing(color)
		newPieceMap[NewPosition(6, castleRank)] = NewRook(color)

	} else if stringAN == "0-0-0" {
		// QUEENSIDE CASTLING:
		castleRank := 0
		if color == White {
			castleRank = 1
		} else {
			castleRank = 8
		}
		newPieceMap[NewPosition(5, castleRank)] = nil
		newPieceMap[NewPosition(1, castleRank)] = nil
		newPieceMap[NewPosition(3, castleRank)] = NewKing(color)
		newPieceMap[NewPosition(4, castleRank)] = NewRook(color)

	} else {
		pieceType := stringAN[:1]            // P/N/B/R/Q/K
		stringAN = stringAN[1:len(stringAN)] //Remove origPiece char
		origFile := stringAN[:1]
		stringAN = stringAN[1:len(stringAN)] //Remove origFile char
		origRank := stringAN[:1]
		stringAN = stringAN[1:len(stringAN)] //Remove origRank char
		moveType := stringAN[:1]             // - (no capture)/x (capture)/ = (pawn promotion)
		stringAN = stringAN[1:len(stringAN)] //Remove moveType char

		origFileInt := fileToInt(origFile)
		origRankInt, _ := strconv.Atoi(origRank)
		origPosition := NewPosition(origFileInt, origRankInt)

		//piece is moving
		nextFile := stringAN[:1]
		stringAN = stringAN[1:len(stringAN)] //Remove nextFile char
		nextRank := stringAN[:1]
		stringAN = stringAN[1:len(stringAN)] //Remove nextRank char

		nextFileInt := fileToInt(nextFile)
		nextRankInt, _ := strconv.Atoi(nextRank)
		nextPosition := NewPosition(nextFileInt, nextRankInt)

		//halfmoveClock resets on capture
		if moveType == "x" {
			state.halfmoveClock = 0
		}

		newPieceMap[origPosition] = nil

		if pieceType == "P" {
			newPieceMap[nextPosition] = NewPawn(color)
			state.halfmoveClock = 0 //halfmoveClock resets on pawn advance
			rankDelta := nextRankInt - origRankInt
			if rankDelta == 2 || rankDelta == -2 { //pawn moved out 2, en passant target
				targetPos := origFile
				if origRankInt == 2 {
					targetPos += "3"
				} else {
					targetPos += "6"
				}
				state.enPassantTarget = AlgebraicMove(targetPos)
			}
			if len(stringAN) > 0 {
				//PAWN PROMOTION
				stringAN = stringAN[1:len(stringAN)] //Remove '='
				promoteType := stringAN[:1]          // Q/N/R/B

				if promoteType == "Q" {
					newPieceMap[nextPosition] = NewQueen(color)
				} else if promoteType == "N" {
					newPieceMap[nextPosition] = NewKnight(color)
				} else if promoteType == "R" {
					newPieceMap[nextPosition] = NewRook(color)
				} else if promoteType == "B" {
					newPieceMap[nextPosition] = NewBishop(color)
				}
			}
		} else if pieceType == "N" {
			newPieceMap[nextPosition] = NewKnight(color)
		} else if pieceType == "B" {
			newPieceMap[nextPosition] = NewBishop(color)
		} else if pieceType == "R" {
			newPieceMap[nextPosition] = NewRook(color)
		} else if pieceType == "Q" {
			newPieceMap[nextPosition] = NewQueen(color)
		} else if pieceType == "K" {
			newPieceMap[nextPosition] = NewKing(color)
		}

	} //esle

	//activeColor change
	if color == White {
		state.activeColor = Black
	} else {
		state.activeColor = White
		state.fullmoveNumber += 1 //fullmove number increments after black turn
	}

	newFEN := state.ConvertToFEN() //convert back to FEN for return
	return newFEN
}

func fileToInt(file string) int {
	if file == "a" || file == "A" {
		return 1
	} else if file == "b" || file == "B" {
		return 2
	} else if file == "c" || file == "C" {
		return 3
	} else if file == "d" || file == "D" {
		return 4
	} else if file == "e" || file == "E" {
		return 5
	} else if file == "f" || file == "F" {
		return 6
	} else if file == "g" || file == "G" {
		return 7
	} else if file == "h" || file == "H" {
		return 8
	}
	return 0
} //fileToInt

func (s *GameState) ValidMoves(pos Position) []Position {
	validMoves := []Position{} //make empty array of positions
	piece := s.PieceAtPosition(pos)
	if piece == nil { //no piece at position
		return []Position{}
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

	//return nil //should not hit this
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
