package queries

import (
	"foodtastechess/game"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	GetGameStateManager() game.GameStateManager
}

type SystemQueryService struct {
	gameStateManager game.GameStateManager `inject:"gameStateManager"`
}

func NewSystemQueryService() *SystemQueryService {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) AnswerQuery(query Query) interface{} {
	return nil
}

func (s *SystemQueryService) GetGameStateManager() game.GameStateManager {
	return s.gameStateManager
}
