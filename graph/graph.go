package graph

import (
	"github.com/facebookgo/inject"
)

type Graph interface {
	Add(name string, object Object) error
	Populate() error
}

type Object interface {
	PreInit(provide Provider) error
	Init() error
}

type Provider func(name string, object Object) error

type injectGraph struct {
	graph inject.Graph
}

func New() Graph {
	return new(injectGraph)
}

func (g *injectGraph) Add(name string, object Object) error {
	if err := object.PreInit(g.Add); err != nil {
		return err
	}

	if err := g.graph.Provide(&inject.Object{
		Name:  name,
		Value: object,
	}); err != nil {
		return err
	}

	return nil
}

func (g *injectGraph) Populate() error {
	err := g.graph.Populate()
	return err
}
