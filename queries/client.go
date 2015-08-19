package queries

import (
	"foodtastechess/game"
	"foodtastechess/users"
)

// ClientQueries is the interface by which other parts of the system
// may get information about current and past states of games.
type ClientQueries interface {
	GameInformation(id game.Id) GameInformation
}

// ClientQueryService provides a concrete implementation of the
// ClientQueries interface.
//
// It composes an injected SystemQueries interface that it uses to
// aggregate information for the supported methods
type ClientQueryService struct {
	Users         users.Users   `inject:"users"`
	SystemQueries SystemQueries `inject:"systemQueries"`
}

func NewClientQueryService() *ClientQueryService {
	cqs := new(ClientQueryService)
	return cqs
}

// GameInformation is a structural representation of the current
// state of a game.
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
	gamePlayers := s.SystemQueries.AnswerQuery(gamePlayersQ).(map[game.Color]string)

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

type MoveRecord struct {
	Move                game.AlgebraicMove
	ResultingBoardState game.FEN
}

func (s *ClientQueryService) GameHistory(id game.Id) []MoveRecord {
	return []MoveRecord{}
}
