package queries

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"foodtastechess/config"
	"foodtastechess/directory"
	"foodtastechess/logger"
)

type QueriesCacheTestSuite struct {
	suite.Suite

	log   *logging.Logger
	cache *queriesCache
}

func (suite *QueriesCacheTestSuite) SetupTest() {
	suite.log = logger.Log("cache_test")

	var (
		d     directory.Directory
		cache *queriesCache
	)

	cache = NewQueriesCache().(*queriesCache)

	d = directory.New()
	d.AddService("configProvider", config.NewConfigProvider("testconfig", "../"))
	d.AddService("queriesCache", cache)

	if err := d.Start(); err != nil {
		suite.log.Fatalf("Could not start directory: %v", err)
	}

	suite.cache = cache
	suite.cache.collection.Remove(map[string]string{})
}

func (suite *QueriesCacheTestSuite) TestStore() {
	var (
		result string = "result!"

		query   *testQuery = newTestQuery("store")
		partial *testQuery = newTestQuery("store")
	)

	query.setResult(result)
	suite.cache.Store(query)

	found := suite.cache.Get(partial)

	assert := assert.New(suite.T())
	assert.Equal(true, found)
	assert.Equal(query.Result, partial.Result)
}

func (suite *QueriesCacheTestSuite) TestGetOrder() {
	var (
		query1 *testQuery = newTestQuery("order")
		query2 *testQuery = newTestQuery("order")
		query3 *testQuery = newTestQuery("order")

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

	partial := newTestQuery("order")

	found := suite.cache.Get(partial)

	assert := assert.New(suite.T())
	assert.Equal(true, found)
	assert.Equal(query2.ComputedAt, partial.ComputedAt)
}

func (suite *QueriesCacheTestSuite) TestDelete() {
	var (
		query *testQuery = newTestQuery("delete")
	)

	suite.cache.Store(query)
	suite.cache.Delete(query)
	found := suite.cache.Get(query)
	assert := assert.New(suite.T())
	assert.Equal(false, found)
}

// testQuery is a test struct that implements Query, so we can use it
// as a pretend query.
type testQuery struct {
	queryRecord `bson:",inline"`

	Param interface{}

	Answered bool
	Result   interface{}
}

// newTestQuery creates a new testQuery with some parameters.
// if only 1 param is provided, testQuery.param will be equal
// to that parameter.
// otherwise, testQuery's param will be set to a slice of
// params.
func newTestQuery(params ...interface{}) *testQuery {
	query := new(testQuery)

	switch len(params) {
	case 0:
		query.Param = nil
	case 1:
		query.Param = params[0]
	default:
		query.Param = params
	}

	return query
}

func (q *testQuery) setResult(result interface{}) {
	q.Result = result
	q.Answered = true
}

func (q *testQuery) deleteResult() {
	q.Result = nil
	q.Answered = false
}

func (q *testQuery) hash() string {
	return fmt.Sprintf("test-query-%v", q.Param)
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

// Entrypoint
func TestQueriesCache(t *testing.T) {
	suite.Run(t, new(QueriesCacheTestSuite))
}
