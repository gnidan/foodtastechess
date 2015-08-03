package main

import (
	"github.com/facebookgo/inject"
	"github.com/stretchr/testify/assert"
	"testing"

	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server"
)

type TestService struct {
	HttpServer         *server.Server              `inject:"httpServer"`
	ClientQueryService *queries.ClientQueryService `inject:"clientQueries"`
	SystemQueryService *queries.SystemQueryService `inject:"systemQueries"`
}

// TestServices sets up a TestService struct that gets provided
// to the graph in order to be injected with the various
// services we expect to be initialized by the app
func TestServices(t *testing.T) {
	var service TestService

	log = logger.Log("main_test")

	app = new(App)

	if err := app.graph.Provide(
		&inject.Object{Value: &service},
	); err != nil {
		log.Fatalf("Could not provide values (%v)", err)
	}

	app.initServices()

	assert.NotNil(t, service.HttpServer)
	assert.NotNil(t, service.ClientQueryService)
	assert.NotNil(t, service.SystemQueryService)
}
