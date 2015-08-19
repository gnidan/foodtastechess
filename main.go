package main

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
	"os/signal"

	"foodtastechess/config"
	"foodtastechess/directory"
	"foodtastechess/events"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server"
	"foodtastechess/users"
)

var (
	app *App
	log *logging.Logger
)

type App struct {
	config    config.ConfigProvider
	directory directory.Directory
	StopChan  chan bool `inject:"stopChan"`
}

func NewApp() *App {
	wd, _ := os.Getwd()

	app := new(App)
	app.directory = directory.New()
	app.config = config.NewConfigProvider("config", wd)
	return app
}

func (app *App) LoadServices() {
	var err error

	app.StopChan = make(chan bool, 1)

	systemQueries := queries.NewSystemQueryService().(*queries.SystemQueryService)
	systemQueries.Complete = true

	services := map[string](interface{}){
		"configProvider":  app.config,
		"httpServer":      server.New(),
		"clientQueries":   queries.NewClientQueryService(),
		"systemQueries":   systemQueries,
		"eventSubscriber": queries.NewQueryBuffer(),
		"users":           users.NewUsers(),
		"events":          events.NewEvents(),
		"stopChan":        app.StopChan,
	}

	for name, value := range services {
		err = app.directory.AddService(name, value)
		if err != nil {
			msg := fmt.Sprintf("Adding %s service failed: %v", name, err)
			log.Error(msg)
		}
	}

	err = app.directory.Start()
	if err != nil {
		msg := fmt.Sprintf("Could not start directory: %v", err)
		log.Error(msg)
		return
	}
}

func (app *App) Start() {
	err := app.directory.Start("httpServer")
	if err != nil {
		msg := fmt.Sprintf("Could not start HTTP Server: %v", err)
		log.Error(msg)
		return
	}

	err = app.directory.Start("eventSubscriber")
	if err != nil {
		msg := fmt.Sprintf("Could not start event subscriber: %v", err)
		log.Error(msg)
		return
	}
}

func (app *App) Stop() {
	err := app.directory.Stop("httpServer")
	if err != nil {
		msg := fmt.Sprintf("Could not stop HTTP Server: %v", err)
		log.Error(msg)
		return
	}

	err = app.directory.Stop("eventSubscriber")
	if err != nil {
		msg := fmt.Sprintf("Could not stop event subscriber: %v", err)
		log.Error(msg)
		return
	}
}

func main() {
	log = logger.Log("main")

	app = NewApp()

	log.Notice("Loading Services")
	app.LoadServices()

	log.Notice("Starting foodtastechess")
	app.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-app.StopChan:
			log.Notice("Quitting.")
			return
		case <-c:
			fmt.Println("")
			app.Stop()
		}
	}
}
