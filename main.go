package main

import (
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"

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
	graph      inject.Graph
	HttpServer *server.Server `inject:"httpServer"`
}

func (app *App) initServices() {
	app.addDependency("app", app)
	app.addDependency("httpServer", server.New())
	app.addDependency("clientQueries", queries.NewClientQueryService())
	app.addDependency("systemQueries", queries.NewSystemQueryService())

	if err := app.graph.Populate(); err != nil {
		log.Error("Could not populate graph (%v)", err)
	}
}

// addDependency is a private interface by which the application
// can add services to its graph for injection.
//
// Each dependency requires a name and of course the service
// itself.
func (app *App) addDependency(name string, value interface{}) {
	object := inject.Object{
		Name:  name,
		Value: value,
	}

	if err := app.graph.Provide(&object); err != nil {
		log.Fatalf(
			"Could not provide value for name %s, got err: %v",
			name,
			err,
		)
	}
}

func main() {
	app = new(App)

	readConf()

	log = logger.Log("main")
	log.Notice("Starting foodtastechess")

	log.Info("Initializing services")
	app.initServices()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	app.HttpServer.Serve("0.0.0.0", port)
}
