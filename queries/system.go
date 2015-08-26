package queries

import (
	"github.com/op/go-logging"
	"reflect"
	"time"

	"foodtastechess/directory"
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/logger"
)

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
	computeAnswer(query Query, skipSearch bool)
	getDependentQueryLookup(query Query) QueryLookup
	getGameCalculator() game.GameCalculator
	getEvents() events.Events
}

type SystemQueryService struct {
	log            *logging.Logger
	GameCalculator game.GameCalculator `inject:"gameCalculator"`
	Events         events.Events       `inject:"events"`
	Cache          Cache               `inject:"queriesCache"`
}

func (s *SystemQueryService) PreProvide(provide directory.Provider) error {
	return provide("queriesCache", NewQueriesCache())
}

func NewSystemQueryService() SystemQueries {
	sqs := new(SystemQueryService)
	sqs.log = logger.Log("systemqueries")
	return sqs
}

func (s *SystemQueryService) AnswerQuery(query Query) interface{} {
	found := s.Cache.Get(query)
	if found {
		s.log.Info("Query %s retrieved from cache: %v", query.hash(), query.getResult())
	} else {
		s.computeAnswer(query, true)
		s.log.Info("Query %s computed: %v", query.hash(), query.getResult())
	}
	return query.getResult()
}

func (s *SystemQueryService) computeAnswer(query Query, skipSearch bool) {
	if !skipSearch && s.Cache.Get(query) {
		s.log.Info("Query %s invalidated, recomputing", query.hash())
		s.Cache.Delete(query)
	}

	query.computeResult(s)

	canMarkComputedAt := reflect.ValueOf(query).
		Elem().
		FieldByName("ComputedAt").
		CanSet()
	if canMarkComputedAt {
		reflect.ValueOf(query).
			Elem().
			FieldByName("ComputedAt").
			Set(reflect.ValueOf(time.Now()))
	}

	s.Cache.Store(query)
}

func (s *SystemQueryService) getDependentQueryLookup(query Query) QueryLookup {
	dependentQueries := query.getDependentQueries()
	for _, dependentQuery := range dependentQueries {
		s.AnswerQuery(dependentQuery)
	}

	return NewQueryLookup(dependentQueries...)
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
