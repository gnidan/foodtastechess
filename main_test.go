package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server"
)

// TestService is a placeholder struct that requires a bunch
// of injected services.
//
// Used to test that the application is actually initializing
// these services, and that if we provide a new dependent
// component, it will be populated accordingly.
type TestService struct {
	HttpServer         *server.Server              `inject:"httpServer"`
	ClientQueryService *queries.ClientQueryService `inject:"clientQueries"`
	SystemQueryService *queries.SystemQueryService `inject:"systemQueries"`
}

// TestServices sets up a TestService struct that gets provided
// to the graph in order to be injected with the various
// services we expect to be initialized by the app
func TestServices(t *testing.T) {
	log = logger.Log("main_test")

	app = newApp()

	var service TestService
	app.graph.Add("testService", &service)

	app.Init()

	assert := assert.New(t)

	assert.NotNil(service.HttpServer)
	assert.NotNil(service.ClientQueryService)
	assert.NotNil(service.SystemQueryService)
}
