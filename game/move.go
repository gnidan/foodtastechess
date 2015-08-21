package game

import (
	"reflect"
	"strconv"
)

//BoundMove - moves max of one position offset
//UnboundMove - moves until hit own piece, border, or opponent + 1
//AdvancingMove - Unique to Pawns: like bound move, but no capture allowed
//CapturingMove - Unique to Pawns: single capturing diagonal move forward
//EnPassantMove - Unique to Pawns: en passant
//CastlingMove - Unique to Rook, King (need for others???) - Castling

type Move interface {
	Translate(pos Position, s *GameState) []AlgebraicMove
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
func AfterMove(algebraicMove AlgebraicMove, prevFEN FEN) FEN {
	state := prevFEN.ConvertToState()
	state.enPassantTarget = ""
	state.halfmoveClock += 1
	color := state.activeColor
	newPieceMap := state.pieceMap
	stringAN := string(algebraicMove)

	if stringAN == "0-0" {
		// KINGSIDE CASTLING:
		castleRank := 0
		updatedCastlingOptions := []string{}
		if color == White {
			castleRank = 1
			//update castling availability
			for _, castleOption := range state.castlingAvailability {
				if castleOption == "k" || castleOption == "q" {
					updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
				}
			} //rof
		} else {
			castleRank = 8
			//update castling availability
			for _, castleOption := range state.castlingAvailability {
				if castleOption == "K" || castleOption == "Q" {
					updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
				}
			} //rof
		}
		state.castlingAvailability = updatedCastlingOptions
		newPieceMap[NewPosition(5, castleRank)] = nil
		newPieceMap[NewPosition(8, castleRank)] = nil
		newPieceMap[NewPosition(7, castleRank)] = NewKing(color)
		newPieceMap[NewPosition(6, castleRank)] = NewRook(color)

	} else if stringAN == "0-0-0" {
		// QUEENSIDE CASTLING:
		castleRank := 0
		updatedCastlingOptions := []string{}
		if color == White {
			castleRank = 1
			//update castling availability
			for _, castleOption := range state.castlingAvailability {
				if castleOption == "k" || castleOption == "q" {
					updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
				}
			} //rof
		} else {
			castleRank = 8
			//update castling availability
			for _, castleOption := range state.castlingAvailability {
				if castleOption == "K" || castleOption == "Q" {
					updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
				}
			} //rof
		}
		state.castlingAvailability = updatedCastlingOptions
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

		//remove moved piece from original position
		newPieceMap[origPosition] = nil

		//update castling availability (if necessary) if King or Rook moved
		if pieceType == "K" || pieceType == "R" {
			updatedCastlingOptions := []string{}
			if color == White {
				for _, castleOption := range state.castlingAvailability {
					if castleOption == "k" || castleOption == "q" {
						updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
					}
				}
			} else if color == Black {
				for _, castleOption := range state.castlingAvailability {
					if castleOption == "K" || castleOption == "Q" {
						updatedCastlingOptions = append(updatedCastlingOptions, castleOption)
					}
				}
			}
			state.castlingAvailability = updatedCastlingOptions
		} //fi

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
				state.enPassantTarget = targetPos
			}
			if len(stringAN) > 0 {
				//pawn promotion/en passant
				nextChar := stringAN[:1]
				stringAN = stringAN[1:len(stringAN)] //Remove '='

				if nextChar == "=" {
					//PAWN PROMOTION
					promoteType := stringAN[:1] // Q/N/R/B
					if promoteType == "Q" {
						newPieceMap[nextPosition] = NewQueen(color)
					} else if promoteType == "N" {
						newPieceMap[nextPosition] = NewKnight(color)
					} else if promoteType == "R" {
						newPieceMap[nextPosition] = NewRook(color)
					} else if promoteType == "B" {
						newPieceMap[nextPosition] = NewBishop(color)
					}
				} else if nextChar == "." {
					//EN PASSANT
					//remove captured pawn
					if color == White {
						capturePosition := NewPosition(nextFileInt, nextRankInt-1)
						newPieceMap[capturePosition] = nil
					} else if color == Black {
						capturePosition := NewPosition(nextFileInt, nextRankInt+1)
						newPieceMap[capturePosition] = nil
					}
				} //else if
			} //fi
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

func intToFile(file int) string {
	if file == 1 {
		return "a"
	} else if file == 2 {
		return "b"
	} else if file == 3 {
		return "c"
	} else if file == 4 {
		return "d"
	} else if file == 5 {
		return "e"
	} else if file == 6 {
		return "f"
	} else if file == 7 {
		return "g"
	} else if file == 8 {
		return "h"
	}
	return ""
} //intToFile

func AllValidMoves(fen FEN) []AlgebraicMove {
	state := fen.ConvertToState()

	fullMovesList := []AlgebraicMove{}
	for file := 1; file <= 8; file++ {
		for rank := 1; rank <= 8; rank++ {
			movesForPosition := state.ValidMovesAtPos(NewPosition(file, rank))
			if len(movesForPosition) > 0 {
				log.Debug("Moves for %v: %v", NewPosition(file, rank), movesForPosition)
			}
			fullMovesList = append(fullMovesList, movesForPosition...)
		} //for rank
	} //for file

	log.Debug("moves list: %v", fullMovesList)

	return fullMovesList
}

//all possible moves disregarding self-checkmate
func AllPossibleMoves(fen FEN) []AlgebraicMove {
	state := fen.ConvertToState()

	fullMovesList := []AlgebraicMove{}
	for file := 1; file <= 8; file++ {
		for rank := 1; rank <= 8; rank++ {
			fullMovesList = append(fullMovesList, state.PossibleMovesAtPos(NewPosition(file, rank))...)
		} //for rank
	} //for file

	return fullMovesList
}

func (s *GameState) ValidMovesAtPos(pos Position) []AlgebraicMove {
	//possibleMoves := []Position{} //make empty array of positions
	piece := s.PieceAtPosition(pos)
	if piece == nil || s.activeColor != piece.Color() { //empty position or opponent piece, no moves
		return []AlgebraicMove{}
	}
	log.Debug("Generating valid moves for position %v", pos)
	moves := piece.Moves()
	possibleMoves := []AlgebraicMove{} //moves that can be executed, disregarding invalid checkmate suicide
	for _, move := range moves {
		possibleMoves = append(possibleMoves, move.Translate(pos, s)...)
	} //rof

	//remove any moves that would put self into check
	validMoves := []AlgebraicMove{} //valid moves, including self-checkmate test
OuterLoop:
	for _, move := range possibleMoves {
		originalFEN := s.ConvertToFEN()
		afterFEN := AfterMove(move, originalFEN)
		afterState := afterFEN.ConvertToState()

		//get own king pos after hypothetical move
		kingPosFile := 0
		kingPosRank := 0
	FindKingLoop:
		for file := 1; file <= 8; file++ {
		FindKingLoopRank:
			for rank := 1; rank <= 8; rank++ {
				piece := afterState.PieceAtPosition(NewPosition(file, rank))
				if piece == nil {
					continue FindKingLoopRank
				} else {
					pieceName := reflect.TypeOf(piece).Name()
					if pieceName == "King" && piece.Color() == s.activeColor { //use original state here, not tempState (color changes)
						kingPosFile = file
						kingPosRank = rank
						break FindKingLoop
					} //fi
				} //esle
			} //for rank
		} //for file

		oppPossibleMoves := AllPossibleMoves(afterFEN) //all possible next opponent moves after hypothetical move

		for _, oppMove := range oppPossibleMoves {
			stringAN := string(oppMove)
			stringAN = stringAN[4:len(stringAN)] //Remove first 4 chars, not needed
			oppFileStr := stringAN[:1]           // file string
			stringAN = stringAN[1:len(stringAN)] //Remove file char
			oppRankStr := stringAN[:1]           // rank string

			oppFile := fileToInt(oppFileStr)
			oppRank, _ := strconv.Atoi(oppRankStr)

			if oppFile == kingPosFile && oppRank == kingPosRank {
				//opponent can take king on next move, invalid
				continue OuterLoop
			} //fi
		} //rof oppMove
		validMoves = append(validMoves, move) //move is valid if it hits here
	} //rof move

	return validMoves
}

//possible moves disregarding self-checkmate
func (s *GameState) PossibleMovesAtPos(pos Position) []AlgebraicMove {
	possibleMoves := []AlgebraicMove{} //make empty array of positions
	piece := s.PieceAtPosition(pos)
	if piece == nil || s.activeColor != piece.Color() { //no piece at position
		return []AlgebraicMove{}
	}
	moves := piece.Moves()
	for _, move := range moves {
		possibleMoves = append(possibleMoves, move.Translate(pos, s)...)
	}
	return possibleMoves
}

func (m *BoundMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	//make sure new position is within boundaries
	if newPos.file < 1 || newPos.file > 8 || newPos.rank < 1 || newPos.rank > 8 {
		return nil
	}

	piece := s.PieceAtPosition(pos)
	pieceName := reflect.TypeOf(piece).Name()
	pieceType := ""
	if pieceName == "Knight" {
		pieceType = "N"
	} else {
		pieceType = pieceName[:1]
	}
	strOrigFile := intToFile(pos.file)
	strOrigRank := strconv.Itoa(pos.rank)
	strNewFile := intToFile(newPos.file)
	strNewRank := strconv.Itoa(newPos.rank)
	moveType := ""
	if s.PieceAtPosition(newPos) == nil {
		//no piece at next position
		moveType = "-"
	} else if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		//opponent piece at new position, capture
		moveType = "x"
	} else {
		return nil
	}
	if pieceType == "P" && (newPos.rank == 1 || newPos.rank == 8) {
		//Pawn promotion, return 4 moves
		algebraicList := []AlgebraicMove{
			AlgebraicMove(pieceType + strOrigFile + strOrigRank + moveType + strNewFile + strNewRank + "=Q"),
			AlgebraicMove(pieceType + strOrigFile + strOrigRank + moveType + strNewFile + strNewRank + "=N"),
			AlgebraicMove(pieceType + strOrigFile + strOrigRank + moveType + strNewFile + strNewRank + "=R"),
			AlgebraicMove(pieceType + strOrigFile + strOrigRank + moveType + strNewFile + strNewRank + "=B"),
		}
		return algebraicList
	} else {
		//not pawn promotion, only return one bound move
		algebraicList := []AlgebraicMove{
			AlgebraicMove(pieceType + strOrigFile + strOrigRank + moveType + strNewFile + strNewRank),
		}
		return algebraicList
	}
} //Translate - BoundMove

func (m *UnboundMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	algebraicList := []AlgebraicMove{}
	nextPos := pos

	piece := s.PieceAtPosition(pos)
	pieceName := reflect.TypeOf(piece).Name()
	pieceType := ""
	if pieceName == "Knight" {
		pieceType = "N"
	} else {
		pieceType = pieceName[:1]
	}
	for {
		nextPos.file = nextPos.file + m.fileOffset
		nextPos.rank = nextPos.rank + m.rankOffset

		strOrigFile := intToFile(pos.file)
		strOrigRank := strconv.Itoa(pos.rank)
		strNewFile := intToFile(nextPos.file)
		strNewRank := strconv.Itoa(nextPos.rank)

		//make sure position is within boundaries
		if nextPos.file < 1 || nextPos.file > 8 || nextPos.rank < 1 || nextPos.rank > 8 {
			//position out of bounds, end of recursion
			return algebraicList
		} else if s.PieceAtPosition(nextPos) == nil {
			//if no piece at position, it is valid.  keep looping
			algebraicList = append(algebraicList, AlgebraicMove(pieceType+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank))
		} else if s.PieceAtPosition(nextPos).Color() != s.PieceAtPosition(pos).Color() {
			//opponent is in position, it is valid, but end of recursion
			return append(algebraicList, AlgebraicMove(pieceType+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank))
		} else {
			//Piece of same color is in next pos.  end of recursion
			return algebraicList
		}
	}
} //Translate - UnboundMove

func (m *FirstPawnMove) Translate(pos Position, s *GameState) []AlgebraicMove {
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
		strOrigFile := intToFile(pos.file)
		strNewFile := intToFile(newPos.file)
		if pawn.Color() == White && pos.rank == 2 {
			return []AlgebraicMove{AlgebraicMove("P" + strOrigFile + "2" + "-" + strNewFile + "4")}
		}
		if pawn.Color() == Black && pos.rank == 7 {
			return []AlgebraicMove{AlgebraicMove("P" + strOrigFile + "7" + "-" + strNewFile + "5")}
		}
	}
	return nil
} //Translate - FirstPawnMove

func (m *AdvancingMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	newPos := Position{
		file: pos.file,
		rank: pos.rank + m.rankDelta,
	}
	algebraicList := []AlgebraicMove{}
	if s.PieceAtPosition(newPos) == nil {
		strOrigFile := intToFile(pos.file)
		strOrigRank := strconv.Itoa(pos.rank)
		strNewFile := intToFile(newPos.file)
		strNewRank := strconv.Itoa(newPos.rank)
		if newPos.rank == 1 || newPos.rank == 8 {
			//Pawn promotion, return 4 moves
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank+"=Q"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank+"=N"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank+"=R"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank+"=B"))
		} else {
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"-"+strNewFile+strNewRank))
		}
	} //fi

	return algebraicList
}

func (m *CapturingMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	//no opponent in position to be captured, invalid
	if s.PieceAtPosition(newPos) == nil {
		return nil
	}
	algebraicList := []AlgebraicMove{}
	//piece in pos must be opponent
	if s.PieceAtPosition(newPos).Color() != s.PieceAtPosition(pos).Color() {
		strOrigFile := intToFile(pos.file)
		strOrigRank := strconv.Itoa(pos.rank)
		strNewFile := intToFile(newPos.file)
		strNewRank := strconv.Itoa(newPos.rank)
		if newPos.rank == 1 || newPos.rank == 8 {
			//Pawn promotion, return 4 moves
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank+"=Q"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank+"=N"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank+"=R"))
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank+"=B"))
		} else {
			algebraicList = append(algebraicList, AlgebraicMove("P"+strOrigFile+strOrigRank+"x"+strNewFile+strNewRank))
		}
	}
	return algebraicList
} //Translate - CapturingMove

func (m *EnPassantMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	newPos := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank + m.rankOffset,
	}
	oppPosition := Position{
		file: pos.file + m.fileOffset,
		rank: pos.rank,
	}

	//no opponent in place for en passant, or position to move in is occupied.  invalid move
	if s.PieceAtPosition(oppPosition) == nil || s.PieceAtPosition(newPos) != nil {
		return nil
	}

	newFile := intToFile(newPos.file)
	newRank := strconv.Itoa(newPos.rank)
	//Check if space to be moved into is En Passant target square
	if s.enPassantTarget != (newFile + newRank) {
		return nil
	}

	//make sure piece in opponent pos is actually opponent
	if s.PieceAtPosition(oppPosition).Color() != s.PieceAtPosition(pos).Color() {
		//Make sure pawn is in correct rank to make en passant
		pawn := s.PieceAtPosition(pos)
		strOrigFile := intToFile(pos.file)
		strNewFile := intToFile(newPos.file)
		if pawn.Color() == White && pos.rank == 5 {
			return []AlgebraicMove{AlgebraicMove("P" + strOrigFile + "5" + "x" + strNewFile + "6" + ".ep")}
		}
		if pawn.Color() == Black && pos.rank == 4 {
			return []AlgebraicMove{AlgebraicMove("P" + strOrigFile + "4" + "x" + strNewFile + "3" + ".ep")}
		}
	}
	return nil
} //Translate - EnPassantMove

func (m *CastlingMove) Translate(pos Position, s *GameState) []AlgebraicMove {
	algebraicList := []AlgebraicMove{} //need to do this in array form, in case 2 castling moves possible

	for _, castleOption := range s.castlingAvailability {
		king := s.PieceAtPosition(pos)
		//Piece must be in initial rank
		if king.Color() == White && pos.file == 5 && pos.rank == 1 {
			//White
			//make sure nothing in between rook and king
			if castleOption == "Q" && s.PieceAtPosition(NewPosition(2, 1)) == nil && s.PieceAtPosition(NewPosition(3, 1)) == nil && s.PieceAtPosition(NewPosition(4, 1)) == nil {
				//Queenside
				algebraicList = append(algebraicList, AlgebraicMove("0-0-0"))
			}
			if castleOption == "K" && s.PieceAtPosition(NewPosition(6, 1)) == nil && s.PieceAtPosition(NewPosition(7, 1)) == nil {
				//Kingside
				algebraicList = append(algebraicList, AlgebraicMove("0-0"))
			}
		} else if king.Color() == Black && pos.file == 5 && pos.rank == 8 {
			//Black
			//make sure nothing in between rook and king
			if castleOption == "q" && s.PieceAtPosition(NewPosition(2, 8)) == nil && s.PieceAtPosition(NewPosition(3, 8)) == nil && s.PieceAtPosition(NewPosition(4, 8)) == nil {
				//Queenside
				algebraicList = append(algebraicList, AlgebraicMove("0-0-0"))
			}
			if castleOption == "k" && s.PieceAtPosition(NewPosition(6, 8)) == nil && s.PieceAtPosition(NewPosition(7, 8)) == nil {
				//Kingside
				algebraicList = append(algebraicList, AlgebraicMove("0-0"))
			}
		} //else if
	} //rof
	return algebraicList

} //Translate - CastlingMove
