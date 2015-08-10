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

	GameStateManager game.GameStateManager `inject:"gameStateManager"`
}

// ComputeAnswer records the call with Query and returns the pre-configured
// mock answer
func (m *MockSystemQueries) AnswerQuery(query Query) interface{} {
	args := m.Called(query)
	return args.Get(0)
}

func (m *MockSystemQueries) GetGameStateManager() game.GameStateManager {
	return nil
}

func (m *MockSystemQueries) IsComplete() bool {
	return m.complete
}

type MockGameStateManager struct {
	mock.Mock
}
