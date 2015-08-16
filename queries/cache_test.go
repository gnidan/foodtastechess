package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"foodtastechess/directory"
	"foodtastechess/game"
	"foodtastechess/logger"
)

type QueriesCacheTestSuite struct {
	suite.Suite

	log   *logging.Logger
	cache *queriesCache
}

type testQuery struct {
	queryRecord

	Result string
}

func (suite *QueriesCacheTestSuite) SetupTest() {
	suite.log = logger.Log("system_test")

	var (
		d     directory.Directory
		cache *queriesCache
	)

	cache = NewQueriesCache().(*queriesCache)

	d = directory.New()
	d.AddService("queriesCache", cache)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.cache = cache
}

func (suite *QueriesCacheTestSuite) TestStore() {
	var (
		gameId game.Id         = 5
		result game.TurnNumber = 6

		query   *turnNumberQuery
		partial *turnNumberQuery
	)

	query = TurnNumberQuery(gameId).(*turnNumberQuery)
	query.Result = result
	suite.cache.Store(query)

	partial = TurnNumberQuery(gameId).(*turnNumberQuery)
	found := suite.cache.Get(partial)

	assert := assert.New(suite.T())
	assert.Equal(true, found)
	assert.Equal(query.Result, partial.Result)
}

func (suite *QueriesCacheTestSuite) TestRetrieveOrdering() {
}

func TestQueriesCache(t *testing.T) {
	suite.Run(t, new(QueriesCacheTestSuite))
}

func (q *testQuery) hasResult() bool {
	return false
}

func (q *testQuery) getResult() interface{} {
	return nil
}

func (q *testQuery) getDependentQueries() []Query {
	return []Query{}
}

func (q *testQuery) computeResult(sqs SystemQueries) {
}

func (q *testQuery) isExpired(now interface{}) bool {
	return false
}

func (q *testQuery) getExpiration(now interface{}) interface{} {
	return nil
}

func (q *testQuery) hash() string {
	return ""
}
