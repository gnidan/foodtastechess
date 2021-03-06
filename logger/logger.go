package logger

import (
	"os"

	"github.com/op/go-logging"
)

type LoggerConfig struct {
	Noise    string
	Level    int
	Location string
	Levels   map[string]int
}

func Log(module string) *logging.Logger {
	log := logging.MustGetLogger(module)
	return log
}

func InitLog(C LoggerConfig) {
	var location *os.File
	switch C.Location {
	case "stderr":
		location = os.Stderr
	case "stdout":
		location = os.Stdout
	}

	backend := logging.NewLogBackend(location, "", 0)
	logging.SetBackend(backend)

	setLevel("", C.Level)
	for name, level := range C.Levels {
		setLevel(name, level)
	}

	var format = logging.MustStringFormatter(
		"[ %{module:12s} ] %{color}%{time:15:04:05.0000} %{level:8s} ▶  %{color:reset}%{message}")

	logging.SetFormatter(format)
}

func setLevel(name string, val int) {
	var level logging.Level
	switch val {
	case 1:
		level = logging.CRITICAL
	case 2:
		level = logging.ERROR
	case 3:
		level = logging.WARNING
	case 4:
		level = logging.NOTICE
	case 5:
		level = logging.INFO
	case 6:
		level = logging.DEBUG
	}
	logging.SetLevel(level, name)
}
