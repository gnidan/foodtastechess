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

func NewClientQueryService() *ClientQueryService {
	cqs := new(ClientQueryService)
	return cqs
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
	turnNumberQ := TurnNumberQuery(id).(*turnNumberQuery)
	turnNumber := s.SystemQueries.AnswerQuery(turnNumberQ).(game.TurnNumber)

	boardStateQ := BoardAtTurnQuery(id, turnNumber).(*boardStateAtTurnQuery)
	boardState := s.SystemQueries.AnswerQuery(boardStateQ).(game.FEN)

	return GameInformation{
		TurnNumber: turnNumber,
		BoardState: boardState,
	}
}
