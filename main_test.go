package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"foodtastechess/commands"
	"foodtastechess/config"
	"foodtastechess/directory"
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/users"
)

type IntegrationTestSuite struct {
	suite.Suite

	log      *logging.Logger
	Cache    queries.Cache
	Commands commands.Commands
	Queries  queries.ClientQueries
	Events   events.Events

	whiteId users.Id
	blackId users.Id
}

func (suite *IntegrationTestSuite) SetupTest() {
	configProvider := config.NewConfigProvider("testconfig", "./")

	suite.log = logger.Log("integration_test")

	suite.Commands = commands.New()
	suite.Queries = queries.NewClientQueryService()

	systemQueries := queries.NewSystemQueryService().(*queries.SystemQueryService)
	eventsService := events.NewEvents().(*events.EventsService)
	usersService := users.NewUsers().(*users.UsersService)

	d := directory.New()
	d.AddService("configProvider", configProvider)
	d.AddService("gameCalculator", game.NewGameCalculator())
	d.AddService("eventSubscriber", queries.NewQueryBuffer())

	d.AddService("systemQueries", systemQueries)
	d.AddService("events", eventsService)
	d.AddService("users", usersService)

	d.AddService("commands", suite.Commands)
	d.AddService("clientQueries", suite.Queries)

	err := d.Start()
	if err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	err = d.Start("eventSubscriber")
	if err != nil {
		msg := fmt.Sprintf("Could not start event subscriber: %v", err)
		log.Error(msg)
		return
	}

	usersService.ResetTestDB()
	eventsService.ResetTestDB()
	systemQueries.Cache.Flush()

	time.Sleep(1 * time.Second)

	white := users.User{
		Uuid:           users.NewId(),
		Name:           "whitePlayer",
		AuthIdentifier: "whiteAuthId",
	}
	usersService.Save(&white)

	black := users.User{
		Uuid:           users.NewId(),
		Name:           "blackPlayer",
		AuthIdentifier: "blackAuthId",
	}
	usersService.Save(&black)

	suite.log.Info("white UUID: %s", white.Uuid)
	suite.log.Info("black UUID: %s", black.Uuid)
	suite.whiteId = white.Uuid
	suite.blackId = black.Uuid
}

func (suite *IntegrationTestSuite) TestGameFlow() {
	assert := assert.New(suite.T())
	var (
		ok     bool
		msg    string
		gameId game.Id
	)

	// Create Game
	ok, msg = suite.Commands.ExecCommand(
		commands.CreateGame, suite.whiteId, map[string]interface{}{
			"color": game.White,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	userGames := suite.Queries.UserGames(suite.whiteId)
	assert.Equal(1, len(userGames))

	gameId = userGames[0]

	// Join Game
	ok, msg = suite.Commands.ExecCommand(
		commands.JoinGame, suite.blackId, map[string]interface{}{
			"gameId": gameId,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	userGames = suite.Queries.UserGames(suite.blackId)
	assert.Equal(1, len(userGames))

	gameInfo, _ := suite.Queries.GameInformation(gameId)
	assert.Equal(queries.GameStatusStarted, gameInfo.GameStatus)

	// Make Move
	ok, msg = suite.Commands.ExecCommand(
		commands.Move, suite.whiteId, map[string]interface{}{
			"gameId": gameId,
			"move":   game.AlgebraicMove("Pb2-b4"),
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	gameHistory, ok := suite.Queries.GameHistory(gameId)
	assert.Equal(true, ok)

	assert.Equal(2, len(gameHistory))
}

func (suite *IntegrationTestSuite) TestDrawOffer() {
	assert := assert.New(suite.T())
	var (
		ok     bool
		msg    string
		gameId game.Id
	)

	// Create Game
	ok, msg = suite.Commands.ExecCommand(
		commands.CreateGame, suite.whiteId, map[string]interface{}{
			"color": game.White,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	userGames := suite.Queries.UserGames(suite.whiteId)
	assert.Equal(1, len(userGames))

	gameId = userGames[0]

	// Join Game
	ok, msg = suite.Commands.ExecCommand(
		commands.JoinGame, suite.blackId, map[string]interface{}{
			"gameId": gameId,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	userGames = suite.Queries.UserGames(suite.blackId)
	assert.Equal(1, len(userGames))

	gameInfo, _ := suite.Queries.GameInformation(gameId)
	assert.Equal(queries.GameStatusStarted, gameInfo.GameStatus)

	assert.Equal(false, gameInfo.OutstandingDrawOffer)

	// Draw Offer
	ok, msg = suite.Commands.ExecCommand(
		commands.OfferDraw, suite.whiteId, map[string]interface{}{
			"gameId": gameId,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	gameInfo, ok = suite.Queries.GameInformation(gameId)
	assert.Equal(true, ok)

	assert.Equal(true, gameInfo.OutstandingDrawOffer)
	assert.Equal(game.White, gameInfo.DrawOfferer)

	// Draw Offer Reject
	ok, msg = suite.Commands.ExecCommand(
		commands.DrawOfferRespond, suite.blackId, map[string]interface{}{
			"gameId": gameId,
			"accept": false,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	gameInfo, ok = suite.Queries.GameInformation(gameId)
	assert.Equal(true, ok)
	assert.Equal(false, gameInfo.OutstandingDrawOffer)
	assert.Equal(queries.GameStatusStarted, gameInfo.GameStatus)

	// Draw Offer 2
	ok, msg = suite.Commands.ExecCommand(
		commands.OfferDraw, suite.blackId, map[string]interface{}{
			"gameId": gameId,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	gameInfo, ok = suite.Queries.GameInformation(gameId)
	assert.Equal(true, ok)

	assert.Equal(true, gameInfo.OutstandingDrawOffer)
	assert.Equal(game.Black, gameInfo.DrawOfferer)

	// Draw Offer 2 Accept
	ok, msg = suite.Commands.ExecCommand(
		commands.DrawOfferRespond, suite.whiteId, map[string]interface{}{
			"gameId": gameId,
			"accept": true,
		},
	)
	assert.Equal(true, ok, msg)

	time.Sleep(100 * time.Millisecond)

	gameInfo, ok = suite.Queries.GameInformation(gameId)
	assert.Equal(true, ok)
	assert.Equal(false, gameInfo.OutstandingDrawOffer)
	assert.Equal(queries.GameStatusEnded, gameInfo.GameStatus)

	assert.Equal(game.GameEndDraw, gameInfo.GameEndReason)
	assert.Equal(game.NoOne, gameInfo.Winner)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
