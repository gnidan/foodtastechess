package queries

import (
	"foodtastechess/game"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	getDependentQueryLookup(query Query) QueryLookup
	getGameCalculator() game.GameCalculator
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

func (s *SystemQueryService) getGameCalculator() game.GameCalculator {
	return s.gameCalculator
}

func (s *SystemQueryService) getDependentQueryLookup(query Query) QueryLookup {
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
