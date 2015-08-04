package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"

	"foodtastechess/graph"
	"foodtastechess/logger"
	"foodtastechess/queries"
	"foodtastechess/server"
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
	graph      graph.Graph
	HttpServer *server.Server `inject:"httpServer"`
}

func newApp() *App {
	app := new(App)

	g := graph.New()
	g.Add("app", app)

	app.graph = g

	return app
}

func (app *App) PreInit(provide graph.Provider) error {
	services := map[string]graph.Object{
		"httpServer":    server.New(),
		"clientQueries": queries.NewClientQueryService(),
		"systemQueries": queries.NewSystemQueryService(),
	}

	for name, object := range services {
		err := provide(name, object)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Init() error {
	err := app.graph.Populate()
	if err != nil {
		return err
	}

	return app.HttpServer.Init()
}

func main() {
	app = newApp()

	readConf()

	log = logger.Log("main")
	log.Notice("Starting foodtastechess")

	log.Info("Initializing App")
	app.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	app.HttpServer.Serve("0.0.0.0", port)
}
