package queries

import (
	"foodtastechess/game"
)

// ClientQueries is the interface by which other parts of the system
// may get information about current and past states of games.
type ClientQueries interface {
	GameInformation(id game.Id) GameInformation
}

// ClientQueryService provides a concrete implementation of the
// ClientQueries interface.
//
// It composes an injected SystemQueries interface that it uses to
// aggregate information for the supported methods
type ClientQueryService struct {
	SystemQueries SystemQueries `inject:"systemQueries"`
}

// GameInformation is a structural representation of the current
// state of a game.
type GameInformation struct {
	Id         game.Id
	TurnNumber game.TurnNumber
	BoardState game.FEN
}

// GameInformation accepts a game ID and queries the SQS for GameInformation
func (s *ClientQueryService) GameInformation(id game.Id) GameInformation {
	var (
		turnNumberQuery Query           = TurnNumberQuery(id)
		turnNumber      game.TurnNumber = s.SystemQueries.GetAnswer(turnNumberQuery).(game.TurnNumber)

		boardStateQuery Query    = BoardAtTurnQuery(id, turnNumber)
		boardState      game.FEN = s.SystemQueries.GetAnswer(boardStateQuery).(game.FEN)
	)

	return GameInformation{
		TurnNumber: turnNumber,
		BoardState: boardState,
	}
}
