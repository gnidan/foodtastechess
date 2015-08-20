package game

import (
	"database/sql/driver"
)

type Id int

func (u *Id) Scan(value interface{}) error {
	*u = Id(value.(int64))
	return nil
}

func (u Id) Value() (driver.Value, error) {
	return int64(u), nil
}

type TurnNumber int

func (u *TurnNumber) Scan(value interface{}) error {
	*u = TurnNumber(value.(int64))
	return nil
}

func (u TurnNumber) Value() (driver.Value, error) {
	return int64(u), nil
}

type FEN string

func (u *FEN) Scan(value interface{}) error {
	*u = FEN(value.([]byte))
	return nil
}

func (u FEN) Value() (driver.Value, error) {
	return string(u), nil
}

type AlgebraicMove string

func (u *AlgebraicMove) Scan(value interface{}) error {
	*u = AlgebraicMove(value.([]byte))
	return nil
}

func (u AlgebraicMove) Value() (driver.Value, error) {
	return string(u), nil
}

type ValidMove interface {
	Move() AlgebraicMove
	PieceMove() (Position, Position)
	ResultingBoardState() FEN
}

type Position interface {
}

type Color int

const (
	White Color = iota
	Black Color = iota
)

type MoveRecord struct {
	Move                AlgebraicMove
	ResultingBoardState FEN
}
