package fixtures

import (
	"github.com/op/go-logging"
	"math/rand"
	"time"

	"foodtastechess/config"
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/users"
)

var log *logging.Logger = logger.Log("fixtures")

type Fixtures struct {
	Config  config.FixturesConfig `inject:"fixturesConfig"`
	Events  events.Events         `inject:"events"`
	Users   users.Users           `inject:"users"`
	Queries queries.ClientQueries `inject:"clientQueries"`
}

func NewFixtures() *Fixtures {
	return new(Fixtures)
}

func (f *Fixtures) Start() error {
	if !f.Config.Enabled {
		return nil
	}

	log.Info("Resetting events")
	eventsService := f.Events.(*events.EventsService)
	eventsService.ResetDB()

	log.Info("Flushing queries cache")
	queriesService := f.Queries.(*queries.ClientQueryService)
	queriesService.FlushCache()

	newUsers := f.userFixtures()
	for _, user := range newUsers {
		f.Users.Save(&user)
	}

	f.runGameFixtures()

	return nil
}

func (f *Fixtures) Stop() error {
	return nil
}

func (f *Fixtures) userFixtures() []users.User {
	usersService := f.Users.(*users.UsersService)
	existingUsers := usersService.GetAll()

	if len(existingUsers) > 8 {
		return []users.User{}
	}

	fakeUserNames := []string{
		"abe",
		"franky g",
		"gerlinde",
		"harry",
		"pauline",
		"jerry",
		"daisuke",
		"chessrobot9000",
	}

	us := []users.User{}
	for _, name := range fakeUserNames {
		us = append(us, users.User{
			Name:           name,
			Uuid:           users.Id(name),
			AuthIdentifier: name,
		})
	}

	return us
}

func (f *Fixtures) runGameFixtures() {
	var (
		numGames            int     = 16
		nextMoveProbability float64 = 0.8
	)

	usersService := f.Users.(*users.UsersService)
	users := usersService.GetAll()

	for gameId := 1; gameId <= numGames; gameId++ {
		shuffle := rand.Perm(len(users))
		white := users[shuffle[0]]
		black := users[shuffle[1]]

		log.Info("Creating game between %s and %s", white.Name, black.Name)

		whiteId := white.Uuid
		blackId := black.Uuid

		if rand.Intn(2) == 0 {
			f.Events.Receive(
				events.NewGameCreateEvent(game.Id(gameId), whiteId, ""),
			)
			time.Sleep(100 * time.Millisecond)
		} else {
			f.Events.Receive(
				events.NewGameCreateEvent(game.Id(gameId), "", blackId),
			)
			time.Sleep(100 * time.Millisecond)
		}

		f.Events.Receive(
			events.NewGameStartEvent(game.Id(gameId), whiteId, blackId),
		)
		time.Sleep(100 * time.Millisecond)

		for turnNumber := 1; ; turnNumber++ {
			if rand.Float64() > nextMoveProbability {
				break
			}

			validMoves, _ := f.Queries.ValidMoves(game.Id(gameId))
			if len(validMoves) == 0 {
				break
			}

			log.Info("Turn %d", turnNumber)

			log.Debug("got valid moves: %v", validMoves)
			randomMoveRecord := validMoves[rand.Intn(len(validMoves))]

			f.Events.Receive(
				events.NewMoveEvent(game.Id(gameId), game.TurnNumber(turnNumber), randomMoveRecord.Move),
			)
			time.Sleep(100 * time.Millisecond)
		}

		log.Info("Done")
	}
}
