package game

import (
	"reflect"
	"strconv"
)

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

func (s GameState) ConvertToFEN() FEN {
	stringFEN := ""

	for rank := 8; rank > 0; rank-- {
		emptyCount := 0
		for file := 1; file <= 8; file++ {

			piece := s.PieceAtPosition(Position{file, rank})

			if piece == nil {
				emptyCount++
				continue
			} else if emptyCount > 0 {
				stringFEN += strconv.Itoa(emptyCount)
				emptyCount = 0
			}
			pieceName := reflect.TypeOf(piece).Name()
			if pieceName == "Pawn" {
				if piece.Color() == White {
					stringFEN += "P"
				} else {
					stringFEN += "p"
				}
			} else if pieceName == "Knight" {
				if piece.Color() == White {
					stringFEN += "N"
				} else {
					stringFEN += "n"
				}
			} else if pieceName == "Bishop" {
				if piece.Color() == White {
					stringFEN += "B"
				} else {
					stringFEN += "b"
				}
			} else if pieceName == "Rook" {
				if piece.Color() == White {
					stringFEN += "R"
				} else {
					stringFEN += "r"
				}
			} else if pieceName == "Queen" {
				if piece.Color() == White {
					stringFEN += "Q"
				} else {
					stringFEN += "q"
				}
			} else if pieceName == "King" {
				if piece.Color() == White {
					stringFEN += "K"
				} else {
					stringFEN += "k"
				}
			}
		}
		if emptyCount > 0 {
			stringFEN += strconv.Itoa(emptyCount)
		}
		if rank > 1 {
			stringFEN += "/"
		}
	}

	newFEN := FEN(stringFEN)
	return newFEN
}

func (fen FEN) ConvertToState() GameState {
	stringFEN := string(fen)
	gameState := GameState{}
	pieceMap := map[Position]Piece{}

	file := 1
	rank := 8

	// 1) PIECE PLACEMENT - from white prespective, starts at rank 8 (top)
	//		Within each rank, each piece is described with it's letter (uppercase for White)
	//		Number of consecutive open spaces in a rank are denoted by a digit (1 through 8)
	//		Each rank is separated by a '/'
	for rank >= 1 {
		next := stringFEN[:1]
		nextInt, _ := strconv.Atoi(next)
		stringFEN = stringFEN[1:len(stringFEN)]

		//Pieces
		if next == "P" {
			pieceMap[NewPosition(file, rank)] = NewPawn(White)
			file += 1
		} else if next == "p" {
			pieceMap[NewPosition(file, rank)] = NewPawn(Black)
			file += 1
		} else if next == "N" {
			pieceMap[NewPosition(file, rank)] = NewKnight(White)
			file += 1
		} else if next == "n" {
			pieceMap[NewPosition(file, rank)] = NewKnight(Black)
			file += 1
		} else if next == "B" {
			pieceMap[NewPosition(file, rank)] = NewBishop(White)
			file += 1
		} else if next == "b" {
			pieceMap[NewPosition(file, rank)] = NewBishop(Black)
			file += 1
		} else if next == "R" {
			pieceMap[NewPosition(file, rank)] = NewRook(White)
			file += 1
		} else if next == "r" {
			pieceMap[NewPosition(file, rank)] = NewRook(Black)
			file += 1
		} else if next == "Q" {
			pieceMap[NewPosition(file, rank)] = NewQueen(White)
			file += 1
		} else if next == "q" {
			pieceMap[NewPosition(file, rank)] = NewQueen(Black)
			file += 1
		} else if next == "K" {
			pieceMap[NewPosition(file, rank)] = NewKing(White)
			file += 1
		} else if next == "k" {
			pieceMap[NewPosition(file, rank)] = NewKing(Black)
			file += 1
		} else if nextInt > 0 && nextInt <= 8 {
			//Digit 1-8, blank spaces
			file += nextInt
		}

		if file > 8 {
			//end of file, go to next rank down
			rank -= 1
			file = 1
		}

	}

	// 2) ACTIVE COLOR -  'w' indicated white moves next, 'b' indicates black

	// 3) CASTLING AVAILABILITY - K/Q/k/q for Black/White King/Queenside; '-' for neither

	// 4) EN PASSANT TARGET SQUARE - in algabraic notation.  If none, then it is '-'.
	//			Recorded as space behind pawn after a double move (even if an opponent isn't near)

	// 5) HALFMOVE CLOCK - number of halfmoves since last capture or pawn advance. Used for 50-move rule

	// 6) FULLMOVE NUMBER - starts at 1, increments after Black's move

	gameState = NewGameState(pieceMap)
	return gameState
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
