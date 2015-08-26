package game

import (
	"database/sql/driver"
	"github.com/op/go-logging"

	"foodtastechess/logger"
)

var log *logging.Logger = logger.Log("game")

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

type Color string

const (
	White Color = "white"
	Black Color = "black"
	NoOne Color = ""
)

func (u *Color) Scan(value interface{}) error {
	*u = Color(value.([]byte))
	return nil
}

func (u Color) Value() (driver.Value, error) {
	return string(u), nil
}

type MoveRecord struct {
	Move                AlgebraicMove
	ResultingBoardState FEN
}

type GameEndReason string

const (
	GameEndConcede   = "concede"
	GameEndDraw      = "stalemate"
	GameEndCheckmate = "checkmate"
)

func (u *GameEndReason) Scan(value interface{}) error {
	*u = GameEndReason(value.([]byte))
	return nil
}

func (u GameEndReason) Value() (driver.Value, error) {
	return string(u), nil
}
