package graph

import (
	"github.com/facebookgo/inject"

	"foodtastechess/common"
)

type Graph interface {
	Add(name string, value interface{}) error
	Populate() error
}

type Object interface {
	PreInit(provide common.Provider) error
	Init() error
}

type injectGraph struct {
	graph inject.Graph
}

func New() Graph {
	return new(injectGraph)
}

func (g *injectGraph) Add(name string, value interface{}) error {
	object, ok := value.(Object)

	if ok {
		if err := object.PreInit(g.Add); err != nil {
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

func (g *injectGraph) Populate() error {
	err := g.graph.Populate()
	return err
}
