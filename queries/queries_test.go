package queries

import (
	"github.com/stretchr/testify/mock"

	"foodtastechess/game"
)

// MockSystemQueries is a mock that we're going to use as a
// SystemQueryInterface
type MockSystemQueries struct {
	mock.Mock

	complete bool

	GameCalculator game.GameCalculator `inject:"gameCalculator"`
}

// ComputeAnswer records the call with Query and returns the pre-configured
// mock answer
func (m *MockSystemQueries) AnswerQuery(query Query) interface{} {
	args := m.Called(query)
	return args.Get(0)
}

func (m *MockSystemQueries) GetComputedDependentQueries(query Query) map[string]Query {
	args := m.Called(query)
	return args.Get(0).(map[string]Query)
}

func (m *MockSystemQueries) GetGameCalculator() game.GameCalculator {
	return m.GameCalculator
}

func (m *MockSystemQueries) IsComplete() bool {
	return m.complete
}

type MockGameCalculator struct {
	mock.Mock
}

func (m *MockGameCalculator) StartingFEN() game.FEN {
	args := m.Called()
	return args.Get(0).(game.FEN)
}

func (m *MockGameCalculator) AfterMove(initial game.FEN, move game.AlgebraicMove) game.FEN {
	args := m.Called(initial, move)
	return args.Get(0).(game.FEN)
}
