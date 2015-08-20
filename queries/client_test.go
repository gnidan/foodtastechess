package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/directory"
	"foodtastechess/game"
	"foodtastechess/logger"
	"foodtastechess/users"
)

// ClientQueriesTestSuite is a collection of tests to ensure the correct
// behavior of the Client Query Service (CQS.) The tests utilize a mocked
// System Query Service (SQS) in order to test that the CQS is merely
// aggregating results from the system.
type ClientQueriesTestSuite struct {
	suite.Suite

	log               *logging.Logger
	mockSystemQueries *MockSystemQueries
	mockUsers         *MockUsers
	clientQueries     ClientQueries
}

// SetupTest prepares the test suite for running by making a fake system
// query service, providing it to a real client query service (the one we
// are testing)
func (suite *ClientQueriesTestSuite) SetupTest() {
	suite.log = logger.Log("client_test")
	var (
		d directory.Directory

		systemQueries MockSystemQueries
		clientQueries ClientQueryService
		mockUsers     MockUsers
	)

	systemQueries.complete = true

	// Set up a directory with:
	//  - A real ClientQueryService (The one we are testing)
	//  - The mocked SystemQueries implementation
	d = directory.New()
	d.AddService("clientQueries", &clientQueries)
	d.AddService("systemQueries", &systemQueries)
	d.AddService("users", &mockUsers)

	// Populate the directory so that clientQueries knows to use our mocked
	// systemQueries
	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory (%v)", err)
	}

	// Store references for use in tests
	suite.mockSystemQueries = &systemQueries
	suite.clientQueries = &clientQueries
	suite.mockUsers = &mockUsers
}

// TestGameInformation tests the ClientQueries.GameInformation() method.
//
// GameInformation should query the SQS for the current turn number and
// the board state at that turn, and return a GameInformation struct
// with that information.
func (suite *ClientQueriesTestSuite) TestGameInformation() {
	var (
		// the game ID we'll be using
		gameId game.Id = 1

		// pretend it's this turn
		expectedTurnNumber game.TurnNumber = 5

		// in this board state
		expectedBoardState game.FEN = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"

		whiteId users.Id = "bob"
		blackId users.Id = "frank"

		gamePlayers map[game.Color]users.Id = map[game.Color]users.Id{
			game.White: whiteId,
			game.Black: blackId,
		}

		expectedWhite users.User = users.User{Uuid: whiteId}
		expectedBlack users.User = users.User{Uuid: blackId}

		// expected query objects we're looking for
		turnNumberQuery  Query = TurnNumberQuery(gameId)
		boardStateQuery  Query = BoardAtTurnQuery(gameId, expectedTurnNumber)
		gamePlayersQuery Query = GamePlayersQuery(gameId)
	)

	// given our expected queries, return our respective expected results
	suite.mockSystemQueries.
		On("AnswerQuery", turnNumberQuery).
		Return(expectedTurnNumber)
	suite.mockSystemQueries.
		On("AnswerQuery", boardStateQuery).
		Return(expectedBoardState)
	suite.mockSystemQueries.
		On("AnswerQuery", gamePlayersQuery).
		Return(gamePlayers)

	suite.mockUsers.
		On("Get", whiteId).
		Return(expectedWhite, true)
	suite.mockUsers.
		On("Get", blackId).
		Return(expectedBlack, true)

	// run the test call
	gameInfo := suite.clientQueries.GameInformation(gameId)

	assert := assert.New(suite.T())
	// and expect that the game info we get back has the pretend values
	assert.Equal(expectedTurnNumber, gameInfo.TurnNumber)
	assert.Equal(expectedBoardState, gameInfo.BoardState)
	assert.Equal(expectedWhite, gameInfo.White)
	assert.Equal(expectedBlack, gameInfo.Black)
}

func (suite *ClientQueriesTestSuite) TestGameHistory() {
	assert := assert.New(suite.T())
	var (
		gameId game.Id = 11

		moves []game.AlgebraicMove = []game.AlgebraicMove{
			"move1",
			"move2",
			"move3",
		}

		states []game.FEN = []game.FEN{
			"start",
			"turn1",
			"turn2",
			"turn3",
		}
	)

	suite.mockSystemQueries.
		On("AnswerQuery", TurnNumberQuery(gameId)).
		Return(game.TurnNumber(len(moves)))

	for i, move := range moves {
		query := MoveAtTurnQuery(gameId, game.TurnNumber(i+1))

		suite.mockSystemQueries.
			On("AnswerQuery", query).
			Return(move)
	}

	for i, state := range states {
		query := BoardAtTurnQuery(gameId, game.TurnNumber(i))

		suite.mockSystemQueries.
			On("AnswerQuery", query).
			Return(state)
	}

	history := suite.clientQueries.GameHistory(gameId)

	assert.Equal(len(states), len(history))

	for i, record := range history {
		var expectedMove game.AlgebraicMove
		if i == 0 {
			expectedMove = ""
		} else {
			expectedMove = moves[i-1]
		}

		expectedState := states[i]

		assert.Equal(expectedMove, record.Move)
		assert.Equal(expectedState, record.ResultingBoardState)
	}
}

func (suite *ClientQueriesTestSuite) TestValidMoves() {
	assert := assert.New(suite.T())
	var (
		gameId game.Id = 13

		turnNumber game.TurnNumber = 39

		validMoves []game.MoveRecord = []game.MoveRecord{
			game.MoveRecord{
				Move:                "finishing move",
				ResultingBoardState: "checkmate",
			},
			game.MoveRecord{
				Move:                "blundering mistake",
				ResultingBoardState: "disaster and famine for years",
			},
		}
	)

	suite.mockSystemQueries.
		On("AnswerQuery", TurnNumberQuery(gameId)).
		Return(turnNumber)

	suite.mockSystemQueries.
		On("AnswerQuery", ValidMovesAtTurnQuery(gameId, turnNumber)).
		Return(validMoves).
		Once()

	result := suite.clientQueries.ValidMoves(gameId)

	assert.Equal(validMoves, result)
}

func TestClientQueriesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientQueriesTestSuite))
}
