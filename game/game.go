package game

import ()

type Id int

type GameInformation struct {
	Id         Id
	TurnNumber TurnNumber
	BoardState FEN
}

type TurnNumber int

type FEN string
