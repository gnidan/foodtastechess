package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"

	"foodtastechess/directory"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server"
	"foodtastechess/user"
)

var (
	app *App
	log *logging.Logger
)

func loggingConf() {
	var C logger.LoggerConfig
	err := viper.MarshalKey("logger", &C)
	if err != nil {
		panic(fmt.Errorf("Can't parse: %s \n", err))
	}
	logger.InitLog(C)
}

func readConf() {
	viper.SetConfigName("config")
	//	viper.AddConfigPath("/etc/appname/") you can use multiple search paths
	viper.AddConfigPath("./")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	loggingConf()
}

type App struct {
	directory directory.Directory
	StopChan  chan bool `inject:"stopChan"`
}

func NewApp() *App {
	app := new(App)
	app.directory = directory.New()
	return app
}

func (app *App) LoadServices() {
	var err error

	app.StopChan = make(chan bool)

	services := map[string](interface{}){
		"httpServer":    server.New(),
		"clientQueries": queries.NewClientQueryService(),
		"systemQueries": queries.NewSystemQueryService(),
		"auth":          user.NewAuthentication(),
		"users":         user.NewUsers(),
		"stopChan":      app.StopChan,
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
}

func main() {
	readConf()

	log = logger.Log("main")

	app = NewApp()

	log.Notice("Loading Services")
	app.LoadServices()

	log.Notice("Starting foodtastechess")
	app.Start()

	<-app.StopChan
}
