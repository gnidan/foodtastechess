package directory

import (
	"errors"
	"fmt"
)

type Directory interface {
	AddService(name string, object interface{}) error
	Start(names ...string) error
	Stop(names ...string) error
}

func New() Directory {
	directory := new(graphDirectory)

	directory.graph = NewGraph()
	directory.services = make(map[string]lifecycleService)
	directory.populated = false

	return directory
}

type graphDirectory struct {
	graph     Graph
	services  map[string]lifecycleService
	populated bool
}

type lifecycleService interface {
	Start() error
	Stop() error
}

func (d *graphDirectory) AddService(name string, service interface{}) error {
	d.populated = false
	err := d.graph.Add(name, service)

	if service, matches := service.(lifecycleService); matches {
		d.services[name] = service
	}

	return err
}

func (d *graphDirectory) Start(names ...string) error {
	if !d.populated {
		err := d.graph.Populate()
		if err != nil {
			return err
		}

		d.populated = true
	}

	for _, name := range names {
		service, ok := d.services[name]
		if !ok {
			msg := fmt.Sprintf("Service '%s' not found", name)
			return errors.New(msg)
		}

		service.Start()
	}

	return nil
}

func (d *graphDirectory) Stop(names ...string) error {
	for _, name := range names {
		service, ok := d.services[name]
		if !ok {
			continue
		}

		err := service.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}
