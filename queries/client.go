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
	turnNumberQuery := Query{
		GameId:    id,
		QueryType: "TurnNumber",
	}
	var turnNumber game.TurnNumber
	turnNumber = game.TurnNumber(s.SystemQueries.GetAnswer(turnNumberQuery).(int))

	return game.GameInformation{
		TurnNumber: turnNumber,
	}
}
