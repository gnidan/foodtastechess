package directory

import (
	"github.com/facebookgo/inject"

	"foodtastechess/common"
)

type graph interface {
	add(name string, value interface{}) error
	populate() error
}

type objectPreProvide interface {
	PreProvide(provide common.Provider) error
}
type objectPostPopulate interface {
	PostPopulate() error
}

type injectGraph struct {
	graph inject.Graph
}

func newGraph() graph {
	return new(injectGraph)
}

func (g *injectGraph) add(name string, value interface{}) error {
	object, ok := value.(objectPreProvide)

	if ok {
		if err := object.PreProvide(g.add); err != nil {
			return err
		}
	}

	if err := g.graph.Provide(&inject.Object{
		Name:  name,
		Value: value,
	}); err != nil {
		return err
	}

	return nil
}

func (g *injectGraph) populate() error {
	if err := g.graph.Populate(); err != nil {
		return err
	}

	for _, injectObject := range g.graph.Objects() {
		value := injectObject.Value

		object, ok := value.(objectPostPopulate)
		if ok {
			err := object.PostPopulate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
