package queries

import (
	"foodtastechess/events"
	"foodtastechess/game"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	getDependentQueryLookup(query Query) QueryLookup
	getGameCalculator() game.GameCalculator
	getEvents() events.Events
}

type SystemQueryService struct {
	GameCalculator game.GameCalculator `inject:"gameCalculator"`
	Events         events.Events       `inject:"eventsService"`
	Cache          Cache               `inject:"queriesCache"`
}

func NewSystemQueryService() SystemQueries {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) AnswerQuery(query Query) interface{} {
	found := s.Cache.Get(query)
	if found {
		return query.getResult()
	}

	query.computeResult(s)
	return query.getResult()
}

func (s *SystemQueryService) getDependentQueryLookup(query Query) QueryLookup {
	return NewQueryLookup()
}

func (s *SystemQueryService) getGameCalculator() game.GameCalculator {
	return s.GameCalculator
}

func (s *SystemQueryService) getEvents() events.Events {
	return s.Events
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
