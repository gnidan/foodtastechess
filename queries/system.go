package queries

import (
	"foodtastechess/game"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	GetDependentQueryLookup(query Query) QueryLookup
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

func (s *SystemQueryService) GetDependentQueryLookup(query Query) QueryLookup {
	return NewQueryLookup()
}

type QueryLookup struct {
	table map[string]Query
}

func NewQueryLookup(queries ...Query) QueryLookup {
	lookup := QueryLookup{}
	lookup.table = make(map[string]Query)

	for _, query := range queries {
		lookup.table[query.hash()] = query
	}

	return lookup
}

func (l QueryLookup) Lookup(query Query) Query {
	return l.table[query.hash()]
}
