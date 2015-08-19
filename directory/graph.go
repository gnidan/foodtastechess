package directory

import (
	"fmt"
	"github.com/facebookgo/inject"
	"sync"
)

type graph interface {
	add(name string, value interface{}) error
	populate() error
}

type objectPreProvide interface {
	PreProvide(provide Provider) error
}
type objectPostPopulate interface {
	PostPopulate() error
}
type objectOverride interface {
	IsComplete() bool
}

type injectGraph struct {
	graph inject.Graph
}

func newGraph() graph {
	return new(injectGraph)
}

func (g *injectGraph) add(name string, value interface{}) error {
	var complete bool

	objectPre, ok := value.(objectPreProvide)
	if ok {
		if err := objectPre.PreProvide(g.add); err != nil {
			return err
		}
	}

	objectOvr, ok := value.(objectOverride)
	if ok {
		complete = objectOvr.IsComplete()
	} else {
		complete = false
	}

	if err := g.graph.Provide(&inject.Object{
		Name:     name,
		Value:    value,
		Complete: complete,
	}); err != nil {
		return err
	}

	return nil
}

func (g *injectGraph) populate() error {
	if err := g.graph.Populate(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, injectObject := range g.graph.Objects() {
		wg.Add(1)

		go func(name string, value interface{}) {
			defer wg.Done()
			object, ok := value.(objectPostPopulate)
			if ok {
				err := object.PostPopulate()
				if err != nil {
					log.Fatal(fmt.Sprintf("Could not perform post-population for %s", name))
				}
			}
		}(injectObject.Name, injectObject.Value)
	}

	wg.Wait()

	return nil
}
