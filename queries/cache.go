package queries

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"reflect"

	"foodtastechess/config"
)

type Cache interface {
	Get(partial Query) bool
	Store(query Query)
	Delete(partial Query)
	Flush()
}

type queriesCache struct {
	Config config.QueriesCacheConfig `inject:"cacheConfig"`

	session    *mgo.Session
	collection *mgo.Collection
}

func (c *queriesCache) PostPopulate() error {
	url := fmt.Sprintf(
		"mongodb://%s:%s",
		c.Config.HostAddr, c.Config.Port,
	)

	session, err := mgo.Dial(url)
	session.SetMode(mgo.Monotonic, true)

	c.session = session
	c.collection = session.DB(c.Config.Database).C("queries")
	return err
}

func NewQueriesCache() Cache {
	return new(queriesCache)
}

func (c *queriesCache) Get(partial Query) bool {
	err := c.collection.
		Find(lookupFor(partial)).
		Sort("-computedat").
		One(partial)
	if err != nil {
		return false
	}

	canMarkAnswered := reflect.ValueOf(partial).Elem().FieldByName("Answered").CanSet()
	if canMarkAnswered {
		reflect.ValueOf(partial).Elem().FieldByName("Answered").SetBool(true)
	}
	return true
}

func (c *queriesCache) Store(query Query) {
	reflect.ValueOf(query).Elem().FieldByName("Hash").SetString(query.hash())

	err := c.collection.Insert(query)
	if err != nil {
		log.Error(fmt.Sprintf("Got error storing: %v", err))
	}
}

func (c *queriesCache) Delete(partial Query) {
	err := c.collection.Remove(lookupFor(partial))
	if err != nil {
		log.Error(
			fmt.Sprintf("Got error deleting %s: %v", partial.hash(), err),
		)
	}
}

func (c *queriesCache) Flush() {
	c.collection.DropCollection()
	c.collection = c.session.DB(c.Config.Database).C("queries")
}

func lookupFor(query Query) map[string]string {
	return map[string]string{
		"hash": query.hash(),
	}
}
