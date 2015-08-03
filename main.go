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
	g   inject.Graph
	a   App
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
	HttpServer *server.Server `inject:"httpServer"`
}

func initServices() {
	httpServer := server.New()
	clientQueryService := queries.NewClientQueryService()
	systemQueryService := queries.NewSystemQueryService()

	if err := g.Provide(
		&inject.Object{Name: "app", Value: &a},
		&inject.Object{Name: "httpServer", Value: httpServer},
		&inject.Object{Name: "clientQueries", Value: clientQueryService},
		&inject.Object{Name: "systemQueries", Value: systemQueryService},
	); err != nil {
		log.Fatalf("Could not provide values (%v)", err)
	}

	if err := g.Populate(); err != nil {
		log.Fatalf("Could not populate graph (%v)", err)
	}
}

func main() {
	readConf()

	log = logger.Log("main")
	log.Notice("Starting foodtastechess")

	log.Info("Initializing services")
	initServices()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	a.HttpServer.Serve("0.0.0.0", port)
}
