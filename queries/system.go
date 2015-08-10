package queries

import (
	"foodtastechess/game"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	GetComputedDependentQueries(query Query) map[string]Query
	GetGameCalculator() game.GameCalculator
}

type SystemQueryService struct {
	gameCalculator game.GameCalculator `inject:"gameCalculator"`
}

func NewSystemQueryService() *SystemQueryService {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) AnswerQuery(query Query) interface{} {
	return nil
}

func (s *SystemQueryService) GetGameCalculator() game.GameCalculator {
	return s.gameCalculator
}

func (s *SystemQueryService) GetComputedDependentQueries(query Query) map[string]Query {
	return make(map[string]Query)
}
