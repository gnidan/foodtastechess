package queries

import (
	"github.com/op/go-logging"

	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/users"
)

// ClientQueries is the interface by which other parts of the system
// may get information about current and past states of games.
type ClientQueries interface {
	UserGames(user users.User) []game.Id
	GameInformation(id game.Id) GameInformation
	GameHistory(id game.Id) []game.MoveRecord
	ValidMoves(id game.Id) []game.MoveRecord
}

// ClientQueryService provides a concrete implementation of the
// ClientQueries interface.
//
// It composes an injected SystemQueries interface that it uses to
// aggregate information for the supported methods
type ClientQueryService struct {
	log           *logging.Logger
	Users         users.Users   `inject:"users"`
	SystemQueries SystemQueries `inject:"systemQueries"`
}

func NewClientQueryService() *ClientQueryService {
	cqs := new(ClientQueryService)
	cqs.log = logger.Log("clientqueries")
	return cqs
}

// UserGames accepts a user and returns a list of game ID's
func (s *ClientQueryService) UserGames(user users.User) []game.Id {
	var (
		games []game.Id = []game.Id{}
	)

	gamesQ := UserGamesQuery(user.Uuid)
	games = s.SystemQueries.AnswerQuery(gamesQ).([]game.Id)
	return games
}

type GameInformation struct {
	Id         game.Id
	TurnNumber game.TurnNumber
	BoardState game.FEN
	White      users.User
	Black      users.User
}

// GameInformation accepts a game ID and queries the SQS for GameInformation
func (s *ClientQueryService) GameInformation(id game.Id) GameInformation {
	gameInfo := new(GameInformation)

	turnNumberQ := TurnNumberQuery(id)
	turnNumber := s.SystemQueries.AnswerQuery(turnNumberQ).(game.TurnNumber)
	gameInfo.TurnNumber = turnNumber

	boardStateQ := BoardAtTurnQuery(id, turnNumber)
	boardState := s.SystemQueries.AnswerQuery(boardStateQ).(game.FEN)
	gameInfo.BoardState = boardState

	gamePlayersQ := GamePlayersQuery(id)
	gamePlayers := s.SystemQueries.AnswerQuery(gamePlayersQ).(map[game.Color]users.Id)

	white, found := s.Users.Get(gamePlayers[game.White])
	if found {
		gameInfo.White = white
	}

	black, found := s.Users.Get(gamePlayers[game.Black])
	if found {
		gameInfo.Black = black
	}

	return *gameInfo
}

func (s *ClientQueryService) GameHistory(gameId game.Id) []game.MoveRecord {
	var (
		history []game.MoveRecord = []game.MoveRecord{}
	)

	turnNumberQ := TurnNumberQuery(gameId)
	turnNumber := s.SystemQueries.AnswerQuery(turnNumberQ).(game.TurnNumber)

	for i := 0; i <= int(turnNumber); i++ {
		stateQ := BoardAtTurnQuery(gameId, game.TurnNumber(i))
		state := s.SystemQueries.AnswerQuery(stateQ).(game.FEN)

		var move game.AlgebraicMove
		if i == 0 {
			move = ""
		} else {
			moveQ := MoveAtTurnQuery(gameId, game.TurnNumber(i))
			move = s.SystemQueries.AnswerQuery(moveQ).(game.AlgebraicMove)
		}

		record := game.MoveRecord{Move: move, ResultingBoardState: state}

		history = append(history, record)
	}

	return history
}

func (s *ClientQueryService) ValidMoves(gameId game.Id) []game.MoveRecord {
	var (
		validMoves []game.MoveRecord = []game.MoveRecord{}
	)
	turnNumberQ := TurnNumberQuery(gameId)
	turnNumber := s.SystemQueries.AnswerQuery(turnNumberQ).(game.TurnNumber)

	validMovesQ := ValidMovesAtTurnQuery(gameId, turnNumber)
	validMoves = s.SystemQueries.AnswerQuery(validMovesQ).([]game.MoveRecord)

	return validMoves
}
