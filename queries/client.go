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
	turnNumberQuery := TurnNumberQuery(id)

	var turnNumber game.TurnNumber
	turnNumber = s.SystemQueries.GetAnswer(turnNumberQuery).(game.TurnNumber)

	return game.GameInformation{
		TurnNumber: turnNumber,
	}
}
