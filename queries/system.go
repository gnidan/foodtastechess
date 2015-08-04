package queries

import (
	"foodtastechess/graph"
)

type SystemQueries interface {
	GetAnswer(query Query) Answer
}

type SystemQueryService struct {
}

func NewSystemQueryService() *SystemQueryService {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) PreInit(provide graph.Provider) error {
	return nil
}

func (s *SystemQueryService) Init() error {
	return nil
}

func (s *SystemQueryService) GetAnswer(query Query) Answer {
	return nil
}
