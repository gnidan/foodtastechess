package main

import (
	"fmt"
	"github.com/facebookgo/inject"
	"github.com/spf13/viper"
	"os"

	"foodtastechess/logger"
	"foodtastechess/server"
)

func loggingConf() {
	var C logger.LoggerConfig
	err := viper.MarshalKey("logger", &C)
	if err != nil {
		panic(fmt.Errorf("Can't parse: %s \n", err))
	}
	logger.InitLog(C)
	log := logger.Log("main")
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
	httpServer server.Server `inject:""`
}

func main() {
	readConf()

	var g inject.Graph
	var a App

	// Here the Populate call is creating instances of NameAPI &
	// PlanetAPI, and setting the HTTPTransport on both to the
	// http.DefaultTransport provided above:
	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	a.httpServer.Serve("0.0.0.0", port)

}
