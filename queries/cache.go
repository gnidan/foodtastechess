package queries

import (
	"fmt"
	"gopkg.in/mgo.v2"

	"foodtastechess/directory"
)

type Cache interface {
	Get(partial Query) bool
	Store(query Query)
	Delete(partial Query)
}

type queriesCache struct {
	Config QueriesCacheConfig `inject:"queriesCacheConfig"`

	session    *mgo.Session
	collection *mgo.Collection
}

func (c *queriesCache) PreProvide(provide directory.Provider) error {
	err := provide("queriesCacheConfig",
		NewMongoDockerComposeConfig(),
	)

	return err
}

func (c *queriesCache) PostPopulate() error {
	url := fmt.Sprintf(
		"mongodb://%s:%s",
		c.Config.HostAddr, c.Config.Port,
	)

	session, err := mgo.Dial(url)
	c.session = session
	c.collection = session.DB(c.Config.Database).C("queries")
	return err
}

func NewQueriesCache() Cache {
	return new(queriesCache)
}

func (c *queriesCache) Get(partial Query) bool {
	return false
}

func (c *queriesCache) Store(query Query) {
}

func (m *queriesCache) Delete(partial Query) {
}
