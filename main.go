package main

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/facebookgo/inject"

	"foodtastechess/logger"
)

func loggingConf() {
	var C logger.LoggerConfig
	err := viper.MarshalKey("logger", &C)
	if err != nil {
		panic(fmt.Errorf("Can't parse: %s \n", err))
	}
	logger.InitLog(C)
	log := logger.Log("main")
	log.Debug("I did it")
	log.Critical("Holy shit")
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

func main() {
	//	s := server.New()

	//	s.Serve("0.0.0.0", os.Getenv("PORT"))
	readConf()
}
