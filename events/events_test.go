package events

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/config"
	"foodtastechess/directory"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/users"
)

type EventsTestSuite struct {
	suite.Suite
	config.ConfigTestProvider

	log            *logging.Logger
	events         *EventsService
	mockSubscriber *MockSubscriber
}

func (suite *EventsTestSuite) SetupTest() {
	suite.InitTestConfig()
	suite.log = logger.Log("events_test")

	var (
		d          directory.Directory
		events     *EventsService
		subscriber *MockSubscriber
	)

	events = NewEvents().(*EventsService)
	subscriber = newMockSubscriber().(*MockSubscriber)

	d = directory.New()
	d.AddService("configProvider", suite.ConfigProvider)
	d.AddService("eventSubscriber", subscriber)
	d.AddService("events", events)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.events = events
	suite.mockSubscriber = subscriber

	suite.events.ResetTestDB()
}

func (suite *EventsTestSuite) TestReceive() {
	var (
		gameId     game.Id            = 5
		turnNumber game.TurnNumber    = 9
		move       game.AlgebraicMove = "Be5"
	)

	event := NewMoveEvent(gameId, turnNumber, move)

	suite.mockSubscriber.
		On("Receive", event).
		Return(nil)

	suite.events.Receive(event)

	suite.mockSubscriber.AssertCalled(suite.T(), "Receive", event)

	events := suite.events.EventsForGame(gameId)

	assert := assert.New(suite.T())
	assert.Equal(1, len(events))
}

func (suite *EventsTestSuite) TestEventsOfTypeForPlayer() {
	var (
		player1 users.Id = "bob"
		player2 users.Id = "frank"

		events []Event = []Event{
			NewGameStartEvent(1, player1, player2),
			NewGameStartEvent(2, player2, player1),
			NewGameStartEvent(3, player2, player1),
			NewGameEndEvent(3, game.GameEndDraw, game.NoOne, player2, player1),
		}
	)

	for _, event := range events {
		suite.mockSubscriber.On("Receive", event).Return(nil).Once()
		suite.events.Receive(event)
	}

	expectedEvents := events[:len(events)-1]
	actualEvents := suite.events.EventsOfTypeForPlayer(player1, GameStartType)

	assert := assert.New(suite.T())
	assert.Equal(len(expectedEvents), len(actualEvents))

	for _, event := range actualEvents {
		assert.Equal(GameStartType, event.Type)
	}
}

// Mocking subscriber
type MockSubscriber struct {
	mock.Mock
}

func newMockSubscriber() EventSubscriber {
	return new(MockSubscriber)
}

func (m *MockSubscriber) Receive(event Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// Entrypoint
func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}
