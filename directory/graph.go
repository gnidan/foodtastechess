package directory

import (
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/mgutz/ansi"
	"sync"
	"time"
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
	for _, injectObject := range g.graph.Objects() {
		if injectObject.Name == name {
			return nil
		}
	}

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
			start := time.Now()
			defer wg.Done()
			object, ok := value.(objectPostPopulate)
			if ok {
				err := object.PostPopulate()
				if err != nil {
					log.Fatal(fmt.Sprintf("Could not perform post-population for %s", name))
				}
			}

			duration := time.Since(start)
			logObjectPopulate(name, duration)
		}(injectObject.Name, injectObject.Value)
	}

	wg.Wait()
	log.Info("Done!")

	return nil
}

var (
	microColor   = ansi.ColorFunc("cyan")
	milliColor   = ansi.ColorFunc("yellow")
	secondsColor = ansi.ColorFunc("red")
)

func logObjectPopulate(name string, duration time.Duration) {
	var (
		value     string
		unit      string
		colorFunc func(string) string

		message string
	)

	s := duration.Seconds()
	ms := duration.Nanoseconds() / 1000000
	us := duration.Nanoseconds() / 1000

	if s >= 1 {
		unit = "s"
		colorFunc = secondsColor
		value = fmt.Sprintf("%.2f", s)
	} else if ms >= 1 {
		unit = "ms"
		colorFunc = milliColor
		value = fmt.Sprintf("%d", ms)
	} else {
		unit = "μs"
		colorFunc = microColor
		value = fmt.Sprintf("%d", us)
	}

	message = fmt.Sprintf(
		"Loaded %s in %s", name, colorFunc(fmt.Sprintf("%s%s", value, unit)),
	)

	log.Info(message)
}
