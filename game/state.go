package game

import (
	"reflect"
	"strconv"
)

type GameState struct {
	pieceMap             map[Position]Piece
	activeColor          Color
	castlingAvailability []string
	enPassantTarget      AlgebraicMove
	halfmoveClock        int
	fullmoveNumber       int
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

	// 1) PIECE PLACEMENT
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

	// 2) ACTIVE COLOR
	stringFEN += " " //add space
	if s.activeColor == White {
		stringFEN += "w"
	} else {
		stringFEN += "b"
	}

	// 3) CASTLING AVAILABILITY
	stringFEN += " " //add space
	if len(s.castlingAvailability) == 0 {
		stringFEN += "-"
	} else {
		for _, castlePiece := range s.castlingAvailability {
			//s.castlingAvailability = []string{"A", castlePiece}
			stringFEN += castlePiece
		}
	} //esle

	// 4) EN PASSANT TARGET SQUARE
	stringFEN += " " //add space
	if s.enPassantTarget == "" {
		stringFEN += "-"
	} else {
		stringFEN += string(s.enPassantTarget)
	}

	// 5) HALFMOVE CLOCK
	stringFEN += " " //add space
	stringFEN += strconv.Itoa(s.halfmoveClock)

	// 6) FULLMOVE NUMBER
	stringFEN += " " //add space
	stringFEN += strconv.Itoa(s.fullmoveNumber)

	newFEN := FEN(stringFEN)
	return newFEN
} //ConvertToFEN

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
LoopPiecePlacement:
	for rank >= 1 {
		next := stringFEN[:1]
		nextInt, _ := strconv.Atoi(next)
		stringFEN = stringFEN[1:len(stringFEN)] //remove char to be processed

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
		if stringFEN[:1] == " " {
			break LoopPiecePlacement
		}
	} //rof
	gameState = NewGameState(pieceMap)

	// 2) ACTIVE COLOR -  'w' indicated white moves next, 'b' indicates black
	stringFEN = stringFEN[1:len(stringFEN)] //Remove space
	if stringFEN[:1] == "w" {
		gameState.activeColor = White
	} else {
		gameState.activeColor = Black
	}
	stringFEN = stringFEN[1:len(stringFEN)] //Remove processed char

	// 3) CASTLING AVAILABILITY - K/Q/k/q for Black/White King/Queenside; '-' for neither
	stringFEN = stringFEN[1:len(stringFEN)] //Remove space
	if stringFEN[:1] == "-" {
		stringFEN = stringFEN[1:len(stringFEN)] //remove dash, don't need to loop
	} else {
	LoopCastlingAvailability:
		for {
			next := stringFEN[:1]
			if next == " " { // end of available castling pieces
				break LoopCastlingAvailability
			}
			gameState.castlingAvailability = append(gameState.castlingAvailability, next)
			stringFEN = stringFEN[1:len(stringFEN)] //Remove processed char
		} //rof
	} //esle

	// 4) EN PASSANT TARGET SQUARE - in algabraic notation.  If none, then it is '-'.
	//			Recorded as space behind pawn after a double move (even if an opponent isn't near)
	stringFEN = stringFEN[1:len(stringFEN)] //Remove space
	next := stringFEN[:1]
	if next == "-" {
		gameState.enPassantTarget = AlgebraicMove("")
		stringFEN = stringFEN[1:len(stringFEN)] //Remove processed char
	} else {
		temp := ""
	LoopEnPassantTarget:
		for {
			next = stringFEN[:1]
			if next != " " {
				temp += next
				stringFEN = stringFEN[1:len(stringFEN)] //Remove processed char
			} else {
				break LoopEnPassantTarget
			}
		} //rof
		gameState.enPassantTarget = AlgebraicMove(temp)
	} //else
	stringFEN = stringFEN[1:len(stringFEN)] //Remove processed char

	// 5) HALFMOVE CLOCK - number of halfmoves since last capture or pawn advance. Used for 50-move rule
	stringFEN = stringFEN[1:len(stringFEN)] //Remove space
	next = stringFEN[:1]
	nextInt, _ := strconv.Atoi(next)
	gameState.halfmoveClock = nextInt

	// 6) FULLMOVE NUMBER - starts at 1, increments after Black's move
	stringFEN = stringFEN[1:len(stringFEN)] //Remove space
	next = stringFEN[:1]
	nextInt, _ = strconv.Atoi(next)
	gameState.fullmoveNumber = nextInt

	return gameState
} //ConvertToState

func InitializeFEN() FEN {
	FEN := FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return FEN
} //InitializeFEN

func InitializeState() GameState {
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

	gameState.activeColor = White
	gameState.castlingAvailability = []string{"K", "Q", "k", "q"}
	gameState.enPassantTarget = AlgebraicMove("")
	gameState.halfmoveClock = 0
	gameState.fullmoveNumber = 1
	return gameState
} //InitializeState
