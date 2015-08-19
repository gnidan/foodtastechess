package directory

import (
	"errors"
	"fmt"

	"foodtastechess/logger"
)

var log = logger.Log("directory")

type Directory interface {
	AddService(name string, object interface{}) error
	Start(names ...string) error
	Stop(names ...string) error
}

func New() Directory {
	directory := new(graphDirectory)

	directory.graph = newGraph()
	directory.services = make(map[string]lifecycleService)
	directory.populated = false

	return directory
}

type Provider func(name string, value interface{}) error

type graphDirectory struct {
	graph     graph
	services  map[string]lifecycleService
	populated bool
}

type lifecycleService interface {
	Start() error
	Stop() error
}

func (d *graphDirectory) AddService(name string, service interface{}) error {
	d.populated = false
	err := d.graph.add(name, service)

	if service, matches := service.(lifecycleService); matches {
		d.services[name] = service
	}

	return err
}

func (d *graphDirectory) Start(names ...string) error {
	if !d.populated {
		err := d.graph.populate()
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
