package game

import ()

type Id int

type TurnNumber int

type FEN string

type AlgebraicMove string

type ValidMoves interface {
}

type Position interface {
}

type Color int

const (
	White Color = iota
	Black Color = iota
)
