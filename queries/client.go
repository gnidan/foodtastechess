package queries

import (
	"foodtastechess/game"
)

type ClientQueries interface {
	GameInformation(id game.Id) GameInformation
}

type ClientQueryService struct {
	SystemQueries SystemQueries `inject:"systemQueries"`
}

type GameInformation struct {
	Id         game.Id
	TurnNumber game.TurnNumber
	BoardState game.FEN
}

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
