package queries

import (
	"foodtastechess/game"
)

type ClientQueries interface {
	GameInformation(id game.Id) game.GameInformation
}

type ClientQueryService struct {
	SystemQueries SystemQueries `inject:"systemQueries"`
}

func (s *ClientQueryService) GameInformation(id game.Id) game.GameInformation {
	var (
		turnNumberQuery Query           = TurnNumberQuery(id)
		turnNumber      game.TurnNumber = s.SystemQueries.GetAnswer(turnNumberQuery).(game.TurnNumber)

		boardStateQuery Query    = BoardAtTurnQuery(id, turnNumber)
		boardState      game.FEN = s.SystemQueries.GetAnswer(boardStateQuery).(game.FEN)
	)

	return game.GameInformation{
		TurnNumber: turnNumber,
		BoardState: boardState,
	}
}
