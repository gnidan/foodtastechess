package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"foodtastechess/config"
	"foodtastechess/directory"
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/users"
)

// Base-class test suite
type QueryTestSuite struct {
	suite.Suite
	config.ConfigTestProvider

	log                *logging.Logger
	mockSystemQueries  *MockSystemQueries
	mockGameCalculator *MockGameCalculator
	mockEvents         *MockEventsService
	mockQueriesCache   *MockQueriesCache
}

func (suite *QueryTestSuite) SetupTest() {
	suite.log = logger.Log("queries_test")

	var (
		d              directory.Directory
		systemQueries  MockSystemQueries
		gameCalculator MockGameCalculator
		events         MockEventsService
		queriesCache   MockQueriesCache
	)

	d = directory.New()
	suite.InitTestConfig()
	d.AddService("configProvider", suite.ConfigProvider)
	d.AddService("systemQueries", &systemQueries)
	d.AddService("gameCalculator", &gameCalculator)
	d.AddService("eventsService", &events)
	d.AddService("queriesCache", &queriesCache)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.mockSystemQueries = &systemQueries
	suite.mockGameCalculator = &gameCalculator
	suite.mockEvents = &events
	suite.mockQueriesCache = &queriesCache
}

// MockUsers is a mock for users service
type MockUsers struct {
	mock.Mock
}

func (m *MockUsers) Get(uuid string) (users.User, bool) {
	args := m.Called(uuid)
	return args.Get(0).(users.User), args.Bool(1)
}

func (m *MockUsers) GetByAuthId(authId string) (users.User, bool) {
	args := m.Called(authId)
	return args.Get(0).(users.User), args.Bool(1)
}

func (m *MockUsers) Save(user *users.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// MockSystemQueries is a mock that we're going to use as a
// SystemQueryInterface
type MockSystemQueries struct {
	mock.Mock

	complete bool

	Events         events.Events       `inject:"eventsService"`
	GameCalculator game.GameCalculator `inject:"gameCalculator"`
	Cache          Cache               `inject:"queriesCache"`
}

// Answer records the call with Query and returns the pre-configured
// mock answer
func (m *MockSystemQueries) AnswerQuery(query Query) interface{} {
	args := m.Called(query)
	return args.Get(0)
}

func (m *MockSystemQueries) computeAnswer(query Query, skipSearch bool) {
	m.Called(query, skipSearch)
}

func (m *MockSystemQueries) getDependentQueryLookup(query Query) QueryLookup {
	args := m.Called(query)
	return args.Get(0).(QueryLookup)
}

func (m *MockSystemQueries) getGameCalculator() game.GameCalculator {
	return m.GameCalculator
}

func (m *MockSystemQueries) getEvents() events.Events {
	return m.Events
}

func (m *MockSystemQueries) IsComplete() bool {
	return m.complete
}

// MockGameCalculator is a mock that is used as a fake
// GameCalculator
type MockGameCalculator struct {
	mock.Mock
}

func (m *MockGameCalculator) StartingFEN() game.FEN {
	args := m.Called()
	return args.Get(0).(game.FEN)
}

func (m *MockGameCalculator) AfterMove(initial game.FEN, move game.AlgebraicMove) game.FEN {
	args := m.Called(initial, move)
	return args.Get(0).(game.FEN)
}

func (m *MockGameCalculator) ValidMoves(state game.FEN) []game.AlgebraicMove {
	args := m.Called(state)
	return args.Get(0).([]game.AlgebraicMove)
}

// MockEventsService is a mock that is used as a fake Events
// service
type MockEventsService struct {
	mock.Mock
}

func (m *MockEventsService) Receive(event events.Event) error {
	return nil
}

func (m *MockEventsService) EventsForGame(gameId game.Id) []events.Event {
	args := m.Called(gameId)
	return args.Get(0).([]events.Event)
}

func (m *MockEventsService) EventsOfTypeForGame(gameId game.Id, eventType events.EventType) []events.Event {
	args := m.Called(gameId, eventType)
	return args.Get(0).([]events.Event)
}

func (m *MockEventsService) EventsOfTypeForPlayer(userId string, eventType events.EventType) []events.Event {
	args := m.Called(userId, eventType)
	return args.Get(0).([]events.Event)
}

func (m *MockEventsService) MoveEventForGameAtTurn(gameId game.Id, turnNumber game.TurnNumber) events.Event {
	args := m.Called(gameId, turnNumber)
	return args.Get(0).(events.Event)
}

// MockQueriesCache is a mock that is used as a fake Cache service
type MockQueriesCache struct {
	mock.Mock
}

func (m *MockQueriesCache) Get(partial Query) bool {
	log.Debug("get %v", partial)
	args := m.Called(partial)
	return args.Bool(0)
}

func (m *MockQueriesCache) Store(query Query) {
	log.Debug("store %v", query)
	m.Called(query)
}

func (m *MockQueriesCache) Delete(partial Query) {
	log.Debug("delete %v", partial)
	m.Called(partial)
}
