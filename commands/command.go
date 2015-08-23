package commands

import (
	"foodtastechess/events"
	"foodtastechess/game"
	"foodtastechess/queries"
)

var (
	commandMap = make(map[string]command)
)

type command struct {
	validators []validator
	gen        eventGenerator
}

type validator func(ctx context, commands Commands) (bool, string)

type eventGenerator func(ctx context, commands Commands) []events.Event

func makeCommand(name string, cmd command) command {
	commandMap[name] = cmd
	return cmd
}

const CreateGame = "create_game"

var createGameCommand = makeCommand(CreateGame, command{
	validators: []validator{},
	gen: func(ctx context, commands Commands) []events.Event {
		gameId := commands.events().NextGameId()

		if ctx.colorChoice == game.White {
			return []events.Event{
				events.NewGameCreateEvent(gameId, ctx.userId, ""),
			}
		} else {
			return []events.Event{
				events.NewGameCreateEvent(gameId, "", ctx.userId),
			}
		}
	},
})

const JoinGame = "join_game"

var joinGameCommand = makeCommand(JoinGame, command{
	validators: []validator{
		gameExists,
		gameNotStarted,
		userNotPlaying,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

		whiteId := gameInfo.White.Uuid
		blackId := gameInfo.Black.Uuid

		if whiteId == "" {
			whiteId = ctx.userId
		} else {
			blackId = ctx.userId
		}

		return []events.Event{
			events.NewGameStartEvent(ctx.gameId, whiteId, blackId),
		}
	},
})

const Move = "move"

var moveCommand = makeCommand(Move, command{
	validators: []validator{
		gameExists,
		userPlaying,
		userActive,
		validMove,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

		return []events.Event{
			events.NewMoveEvent(ctx.gameId, gameInfo.TurnNumber+1, ctx.move),
		}
	},
})

func gameExists(ctx context, commands Commands) (bool, string) {
	_, exists := commands.queries().GameInformation(ctx.gameId)

	if !exists {
		return false, "Game does not exist."
	} else {
		return true, ""
	}
}

func gameDoesNotExist(ctx context, commands Commands) (bool, string) {
	_, exists := commands.queries().GameInformation(ctx.gameId)

	if exists {
		return false, "Game already exists."
	} else {
		return true, ""
	}
}

func gameNotStarted(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
	if gameInfo.GameStatus != queries.GameStatusCreated {
		return false, "Game cannot have been started"
	} else {
		return true, ""
	}
}

func userPlaying(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

	if ctx.userId == gameInfo.White.Uuid {
		return true, ""
	} else if ctx.userId == gameInfo.Black.Uuid {
		return true, ""
	} else {
		return false, "You are not playing in the game."
	}
}

func userNotPlaying(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

	if ctx.userId == gameInfo.White.Uuid || ctx.userId == gameInfo.Black.Uuid {
		return false, "You are already playing this game."
	} else {
		return true, ""
	}
}

func userActive(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

	if gameInfo.ActiveColor == game.White && gameInfo.White.Uuid == ctx.userId {
		return true, ""
	} else if gameInfo.ActiveColor == game.Black && gameInfo.Black.Uuid == ctx.userId {
		return true, ""
	} else {
		return false, "It is not your turn."
	}
}

func validMove(ctx context, commands Commands) (bool, string) {
	validMoves, _ := commands.queries().ValidMoves(ctx.gameId)

	for _, validMove := range validMoves {
		if ctx.move == validMove.Move {
			return true, ""
		}
	}

	return false, "Invalid move"
}
