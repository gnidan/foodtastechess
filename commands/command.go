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

// Commands!

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
		gameNotEnded,
		userPlaying,
		userActive,
		gameHasNoDrawOffer,
		validMove,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
		es := []events.Event{
			events.NewMoveEvent(ctx.gameId, gameInfo.TurnNumber+1, ctx.move),
		}

		whiteId := gameInfo.White.Uuid
		blackId := gameInfo.Black.Uuid

		var player game.Color

		if ctx.userId == whiteId {
			player = game.White
		} else {
			player = game.Black
		}

		var lastChar = string(ctx.move)[len(ctx.move)-1:]

		if lastChar == "#" {
			es = append(es, events.NewGameEndEvent(
				ctx.gameId, game.GameEndCheckmate, player,
				whiteId, blackId,
			))
		} else if lastChar == "S" {
			es = append(es, events.NewGameEndEvent(
				ctx.gameId, game.GameEndDraw, game.NoOne,
				whiteId, blackId,
			))
		}

		return es
	},
})

const Concede = "concede"

var concedeCommand = makeCommand(Concede, command{
	validators: []validator{
		gameExists,
		userPlaying,
		gameStarted,
		gameNotEnded,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

		whiteId := gameInfo.White.Uuid
		blackId := gameInfo.Black.Uuid

		var winner game.Color

		if ctx.userId == whiteId {
			winner = game.Black
		} else {
			winner = game.White
		}

		return []events.Event{
			events.NewGameEndEvent(ctx.gameId, game.GameEndConcede, winner, whiteId, blackId),
		}
	},
})

const OfferDraw = "offer_draw"

var offerDrawCommand = makeCommand(OfferDraw, command{
	validators: []validator{
		gameExists,
		userPlaying,
		gameStarted,
		gameNotEnded,
		gameHasNoDrawOffer,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

		var offerer game.Color
		if ctx.userId == gameInfo.White.Uuid {
			offerer = game.White
		} else {
			offerer = game.Black
		}

		return []events.Event{
			events.NewDrawOfferEvent(ctx.gameId, offerer),
		}
	},
})

const DrawOfferRespond = "respond"

var drawOfferRespondCommand = makeCommand(DrawOfferRespond, command{
	validators: []validator{
		gameExists,
		userPlaying,
		gameStarted,
		gameNotEnded,
		opponentOfferedDraw,
	},

	gen: func(ctx context, commands Commands) []events.Event {
		es := []events.Event{
			events.NewDrawOfferResponseEvent(ctx.gameId, ctx.accept),
		}

		if ctx.accept {
			gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
			gameEnd := events.NewGameEndEvent(
				ctx.gameId, game.GameEndDraw, game.NoOne,
				gameInfo.White.Uuid, gameInfo.Black.Uuid,
			)

			es = append(es, gameEnd)
		}

		return es
	},
})

// Validators!

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

func gameStarted(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
	if gameInfo.GameStatus == queries.GameStatusCreated {
		return false, "Game must have already started."
	} else {
		return true, ""
	}
}

func gameNotStarted(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
	if gameInfo.GameStatus != queries.GameStatusCreated {
		return false, "Game cannot have been started."
	} else {
		return true, ""
	}
}

func gameNotEnded(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)
	if gameInfo.GameStatus == queries.GameStatusEnded {
		return false, "Game cannot have ended."
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

func gameHasNoDrawOffer(ctx context, commands Commands) (bool, string) {
	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

	if gameInfo.OutstandingDrawOffer {
		return false, "There is an outstanding draw offer."
	} else {
		return true, ""
	}
}

func opponentOfferedDraw(ctx context, commands Commands) (bool, string) {
	msg := "Your opponent must have offered draw."

	gameInfo, _ := commands.queries().GameInformation(ctx.gameId)

	if !gameInfo.OutstandingDrawOffer {
		return false, msg
	}

	var userColor game.Color
	if ctx.userId == gameInfo.White.Uuid {
		userColor = game.White
	} else {
		userColor = game.Black
	}

	if gameInfo.DrawOfferer == userColor {
		return false, msg
	} else {
		return true, ""
	}

}
