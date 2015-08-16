package queries

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

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
	queryRecord `bson:",inline"`

	Answered bool
	Result   string
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

func (suite *QueriesCacheTestSuite) TestGetOrder() {
	var (
		query1 *testQuery = new(testQuery)
		query2 *testQuery = new(testQuery)
		query3 *testQuery = new(testQuery)

		time1 time.Time
		time2 time.Time
		time3 time.Time
	)

	time1 = time.Now().Truncate(time.Millisecond)
	time2 = time1.Add(5 * time.Hour).Truncate(time.Millisecond)
	time3 = time1.Add(-3 * time.Hour).Truncate(time.Millisecond)

	query1.ComputedAt = time1
	query2.ComputedAt = time2
	query3.ComputedAt = time3

	suite.cache.Store(query1)
	suite.cache.Store(query2)
	suite.cache.Store(query3)

	partial := new(testQuery)

	found := suite.cache.Get(partial)

	assert := assert.New(suite.T())
	assert.Equal(true, found)
	assert.Equal(query2.ComputedAt, partial.ComputedAt)
}

func TestQueriesCache(t *testing.T) {
	suite.Run(t, new(QueriesCacheTestSuite))
}

func (q *testQuery) hasResult() bool {
	return q.Answered
}

func (q *testQuery) getResult() interface{} {
	return q.Result
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
	return "test-query"
}
